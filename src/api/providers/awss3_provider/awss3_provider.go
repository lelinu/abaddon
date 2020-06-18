package awss3_provider

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lelinu/abaddon/src/api/domain/awss3"
	"github.com/lelinu/abaddon/src/api/utils/file_utils"
	"github.com/lelinu/api_utils/utils/error_utils"
	"os"
	"path/filepath"
	"strings"
)

type AwsS3Provider struct {
	client        *s3.S3
	config        *aws.Config
	isLdapEnabled bool
	endPoint	  string
}

type IAwsS3Provider interface {
	Init(req *awss3.InitRequest) (*AwsS3Provider, *error_utils.ApiError)
	Ls(req *awss3.LsRequest) (*awss3.LsResponse, *error_utils.ApiError)
	Stat(req *awss3.StatRequest) (*awss3.StatResponse, *error_utils.ApiError)
	Cat(req *awss3.CatRequest) (*awss3.CatResponse, *error_utils.ApiError)
	MkDir(req *awss3.MkDirRequest) *error_utils.ApiError
	Rm(req *awss3.RmRequest) *error_utils.ApiError
	Mv(req *awss3.MvRequest) *error_utils.ApiError
	Save(req *awss3.SaveRequest) (*awss3.SaveResponse, *error_utils.ApiError)
	Touch(req *awss3.TouchRequest) *error_utils.ApiError
}

func NewAwsS3Provider(isLdapEnabled bool, endPoint string) IAwsS3Provider {
	return &AwsS3Provider{client: nil, config: nil, isLdapEnabled: isLdapEnabled, endPoint: endPoint}
}

func (s *AwsS3Provider) Init(req *awss3.InitRequest) (*AwsS3Provider, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// if ldap is enabled connect via ldap
	if s.isLdapEnabled{
		req.Username = awss3.NewRgwToken(req.Username, req.Password)
		req.Password = "."
	}

	// build configuration
	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(req.Username, req.Password, ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(req.Region),
		Endpoint:         aws.String(s.endPoint),
	}

	// create new session
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Printf("error is %v", err)
		return nil, error_utils.NewInternalServerError(err.Error())
	}

	client := s3.New(sess)
	s.config = config
	s.client = client

	return s, nil
}

func (s *AwsS3Provider) Ls(req *awss3.LsRequest) (*awss3.LsResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// parse path to aws path
	p := awss3.ParsePath(req.Path)
	files := make([]os.FileInfo, 0)

	if p.Bucket == "" {
		b, err := s.client.ListBuckets(&s3.ListBucketsInput{})
		if err != nil {
			return nil, error_utils.NewInternalServerError(err.Error())
		}

		var canMove = false
		for _, bucket := range b.Buckets {
			files = append(files, &file_utils.File{
				FName:   *bucket.Name,
				FType:   file_utils.DIRECTORY,
				FTime:   bucket.CreationDate.Unix(),
				CanMove: &canMove,
			})
		}
		return awss3.NewLsResponse(files), nil
	}

	session, apiErr := s.createSession(p.Bucket)
	if apiErr != nil {
		return nil, apiErr
	}
	client := s3.New(session)

	objs, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(p.Bucket),
		Prefix:    aws.String(p.Path),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, error_utils.NewInternalServerError(err.Error())
	}

	for i, object := range objs.Contents {
		if i == 0 && *object.Key == p.Path {
			continue
		}

		files = append(files, &file_utils.File{
			FName: filepath.Base(*object.Key),
			FType: file_utils.FILE,
			FTime: object.LastModified.Unix(),
			FSize: *object.Size,
		})
	}
	for _, object := range objs.CommonPrefixes {

		files = append(files, &file_utils.File{
			FName: filepath.Base(*object.Prefix),
			FType: file_utils.DIRECTORY,
		})
	}
	return awss3.NewLsResponse(files), nil
}

func (s *AwsS3Provider) Stat(req *awss3.StatRequest) (*awss3.StatResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// parse path to aws path
	p := awss3.ParsePath(req.Path)
	session, apiErr := s.createSession(p.Bucket)
	if apiErr != nil {
		return nil, apiErr
	}
	client := s3.New(session)

	params := &s3.ListObjectVersionsInput{
		Bucket:    aws.String(p.Bucket),
		Prefix:    aws.String(p.Path),
		Delimiter: aws.String("/"),
	}

	var versions []*s3.ObjectVersion
	err := client.ListObjectVersionsPages(params,
		func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
			for _, item := range page.Versions {
				versions = append([]*s3.ObjectVersion{item}, versions...)
			}
			return !lastPage
		})
	// check if there's an error
	if err != nil {
		return nil, error_utils.NewBadRequestError(err.Error())
	}

	// map s3 object to internal objects
	var objectVersions []*awss3.ObjectVersion
	for _, v := range versions {
		objectVersions = append(objectVersions, awss3.NewObjectVersion(v.ETag, v.IsLatest, v.Key, v.LastModified, v.Size, v.VersionId))
	}

	// return
	return awss3.NewStatResponse(objectVersions), nil
}

func (s *AwsS3Provider) Cat(req *awss3.CatRequest) (*awss3.CatResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// create a new session with new path
	client, p, apiErr := s.createS3ClientAndPath(req.Path)
	if apiErr != nil {
		return nil, apiErr
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(p.Path),
	}

	// if request must use encryption
	if req.UseEncryption() {
		input.SSECustomerAlgorithm = aws.String(req.EncryptionAlgorithm)
		input.SSECustomerKey = aws.String(req.EncryptionKey)
	}

	// call get object
	obj, err := client.GetObject(input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok == false {
			return nil, error_utils.NewInternalServerError(err.Error())
		}

		if awsErr.Code() == awss3.INVALID_REQUEST && strings.Contains(strings.ToLower(awsErr.Message()), "server side encryption") {
			return nil, error_utils.NewBadRequestError("The object was stored using a form of Server Side Encryption and cannot be retrieved")
		} else if awsErr.Code() == awss3.INVALID_ARGUMENT && strings.Contains(awsErr.Message(), "secret key was invalid") {
			return nil, error_utils.NewBadRequestError("This file is encrypted file, you need the correct key")
		} else if awsErr.Code() == awss3.ACCESS_DENIED {
			return nil, error_utils.NewForbiddenError("Access denied")
		} else if awsErr.Code() == awss3.NO_SUCH_KEY {
			return nil, error_utils.NewBadRequestError("File not found")
		}

		return nil, error_utils.NewInternalServerError(err.Error())
	}

	return awss3.NewCatResponse(obj.Body), nil
}

func (s *AwsS3Provider) MkDir(req *awss3.MkDirRequest) *error_utils.ApiError {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return reqErr
	}

	// create client and S3 path
	client, p, apiErr := s.createS3ClientAndPath(req.Path)
	if apiErr != nil {
		return apiErr
	}

	// if path is empty create only a bucket
	if p.Path == "" {
		_, err := client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(req.Path),
		})

		if err != nil {
			return error_utils.NewInternalServerError(err.Error())
		}

		return nil
	}

	// else call put object
	_, err := client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(p.Path),
	})

	return error_utils.NewInternalServerError(err.Error())
}

func (s *AwsS3Provider) Rm(req *awss3.RmRequest) *error_utils.ApiError {

	// create client and S3 path
	client, p, apiErr := s.createS3ClientAndPath(req.Path)
	if apiErr != nil {
		return apiErr
	}

	if p.Bucket == "" {
		return error_utils.NewBadRequestError("Not found")
	}

	objs, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(p.Bucket),
		Prefix:    aws.String(p.Path),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return error_utils.NewInternalServerError(err.Error())
	}
	for _, obj := range objs.Contents {
		_, err := client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(p.Bucket),
			Key:    obj.Key,
		})
		if err != nil {
			return error_utils.NewInternalServerError(err.Error())
		}
	}
	for _, pref := range objs.CommonPrefixes {

		s.Rm(awss3.NewRmRequest("/" + p.Bucket + "/" + *pref.Prefix))

		_, err := client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(p.Bucket),
			Key:    pref.Prefix,
		})
		if err != nil {
			return error_utils.NewInternalServerError(err.Error())
		}
	}
	if err != nil {
		return error_utils.NewInternalServerError(err.Error())
	}

	if p.Path == "" {
		_, err := client.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(p.Bucket),
		})
		if err != nil {
			return error_utils.NewInternalServerError(err.Error())
		}
		return nil
	}
	_, err = client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(p.Path),
	})
	if err != nil {
		return error_utils.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *AwsS3Provider) Mv(req *awss3.MvRequest) *error_utils.ApiError {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return reqErr
	}

	// parse path to aws path
	f := awss3.ParsePath(req.PathFrom)
	t := awss3.ParsePath(req.PathTo)

	session, apiErr := s.createSession(f.Bucket)
	if apiErr != nil {
		return apiErr
	}
	client := s3.New(session)

	// from path must not be empty
	if f.Path == "" {
		return error_utils.NewUnauthorizedError("Can't move this")
	}

	// copy object
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(t.Bucket),
		CopySource: aws.String(f.Bucket + "/" + f.Bucket),
		Key:        aws.String(t.Path),
	}
	if req.UseEncryption() {
		input.CopySourceSSECustomerAlgorithm = aws.String(req.EncryptionAlgorithm)
		input.CopySourceSSECustomerKey = aws.String(req.EncryptionKey)
		input.SSECustomerAlgorithm = aws.String(req.EncryptionAlgorithm)
		input.SSECustomerKey = aws.String(req.EncryptionKey)
	}

	_, err := client.CopyObject(input)
	if err != nil {
		return error_utils.NewInternalServerError(err.Error())
	}
	// remove from object
	return s.Rm(awss3.NewRmRequest(req.PathFrom))
}

func (s *AwsS3Provider) Save(req *awss3.SaveRequest) (*awss3.SaveResponse, *error_utils.ApiError) {
	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// parse path & create session
	p := awss3.ParsePath(req.Path)
	if p.Bucket == "" {
		return nil, error_utils.NewBadRequestError("Can't do that on S3")
	}
	sess, apiErr := s.createSession(req.Path)
	if apiErr != nil {
		return nil, apiErr
	}

	// upload via s3 managers
	uploader := s3manager.NewUploader(sess)
	input := s3manager.UploadInput{
		Body:   req.File,
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(p.Path),
	}
	if req.UseEncryption() {
		fmt.Printf("using encryption %v", req.UseEncryption())
		input.SSECustomerAlgorithm = aws.String(req.EncryptionAlgorithm)
		input.SSECustomerKey = aws.String(req.EncryptionKey)
	}
	resp, err := uploader.Upload(&input)
	if err != nil {
		return nil, error_utils.NewInternalServerError(err.Error())
	}

	return awss3.NewSaveResponse(resp.Location, resp.VersionID, resp.UploadID), nil
}

func (s *AwsS3Provider) Touch(req *awss3.TouchRequest) *error_utils.ApiError {
	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return reqErr
	}

	// parse path & create session
	p := awss3.ParsePath(req.Path)
	if p.Bucket == "" {
		return error_utils.NewBadRequestError("Can't do that on S3")
	}
	sess, apiErr := s.createSession(req.Path)
	if apiErr != nil {
		return apiErr
	}
	client := s3.New(sess)
	input := &s3.PutObjectInput{
		Body:          strings.NewReader(""),
		ContentLength: aws.Int64(0),
		Bucket:        aws.String(p.Bucket),
		Key:           aws.String(p.Path),
	}
	if req.UseEncryption() {
		fmt.Printf("using encryption %v", req.UseEncryption())
		input.SSECustomerAlgorithm = aws.String(req.EncryptionAlgorithm)
		input.SSECustomerKey = aws.String(req.EncryptionKey)
	}
	_, err := client.PutObject(input)
	if err != nil {
		return error_utils.NewInternalServerError(err.Error())
	}
	return nil
}

func (s *AwsS3Provider) createS3ClientAndPath(path string) (*s3.S3, *awss3.Path, *error_utils.ApiError) {
	p := awss3.ParsePath(path)

	session, apiErr := s.createSession(p.Bucket)
	if apiErr != nil {
		return nil, nil, apiErr
	}
	return s3.New(session), p, nil
}

func (s *AwsS3Provider) createSession(bucket string) (*session.Session, *error_utils.ApiError) {

	res, err := s.client.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		fmt.Printf("err is %v", err)
		s.config.Region = aws.String("us-east-1")
	} else {
		if res.LocationConstraint == nil {
			s.config.Region = aws.String("us-east-1")
		} else {
			s.config.Region = res.LocationConstraint
		}
	}

	sess, err := session.NewSession(s.config)
	if err != nil {
		return nil, error_utils.NewInternalServerError(err.Error())
	}
	return sess, nil
}
