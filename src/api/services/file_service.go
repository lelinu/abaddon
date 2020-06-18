package services

import (
	"encoding/base64"
	"fmt"
	"github.com/lelinu/abaddon/src/api/config"
	"github.com/lelinu/abaddon/src/api/domain/awss3"
	"github.com/lelinu/abaddon/src/api/domain/file"
	"github.com/lelinu/abaddon/src/api/middleware"
	"github.com/lelinu/abaddon/src/api/utils/file_utils"
	"github.com/lelinu/api_utils/log/lzap"
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/mime_utils"
	"hash/fnv"
	"os"
	"strconv"
	"time"
)

type fileService struct {
	logger lzap.IService
}

func NewFileService(logger lzap.IService) IFileService {
	return &fileService{logger: logger}
}

type IFileService interface {
	Ls(scope *middleware.Scope, req *file.LsRequest) (*file.LsResponse, *error_utils.ApiError)
	Stat(scope *middleware.Scope, req *file.StatRequest) (*file.StatResponse, *error_utils.ApiError)
	Cat(scope *middleware.Scope, req *file.CatRequest) (*file.CatResponse, *error_utils.ApiError)
	MkDir(scope *middleware.Scope, req *file.MkDirRequest) *error_utils.ApiError
	Rm(scope *middleware.Scope, req *file.RmRequest) *error_utils.ApiError
	Mv(scope *middleware.Scope, req *file.MvRequest) *error_utils.ApiError
	Save(scope *middleware.Scope, req *file.SaveRequest) (*file.SaveResponse, *error_utils.ApiError)
	Touch(scope *middleware.Scope, req *file.TouchRequest) *error_utils.ApiError
}

func (f *fileService) Ls(scope *middleware.Scope, req *file.LsRequest) (*file.LsResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// build request
	request := &awss3.LsRequest{
		Path: req.Path,
	}

	response, apiErr := scope.AwsS3Provider.Ls(request)
	if apiErr != nil {
		return nil, apiErr
	}

	// build response with etag
	entries := response.Files
	if len(entries) == 0 {
		return file.NewLsResponse(nil, ""), nil
	}

	files := make([]file_utils.FileInfo, len(entries))
	etagger := fnv.New32()
	etagger.Write([]byte(req.Path + strconv.Itoa(len(entries))))
	for i := 0; i < len(entries); i++ {
		name := entries[i].Name()
		modTime := entries[i].ModTime().UnixNano() / int64(time.Millisecond)

		if i < 200 {
			etagger.Write([]byte(name + strconv.Itoa(int(modTime))))
		}

		files[i] = file_utils.FileInfo{
			Name: name,
			Size: entries[i].Size(),
			Time: modTime,
			Type: func(mode os.FileMode) string {
				if mode.IsRegular() {
					return file_utils.FILE
				}
				return file_utils.DIRECTORY
			}(entries[i].Mode()),
		}
	}

	etagValue := base64.StdEncoding.EncodeToString(etagger.Sum(nil))
	return file.NewLsResponse(files, etagValue), nil
}

func (f *fileService) Stat(scope *middleware.Scope, req *file.StatRequest) (*file.StatResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// build request
	request := &awss3.StatRequest{
		Path: req.Path,
	}

	// call stat
	response, apiErr := scope.AwsS3Provider.Stat(request)
	if apiErr != nil {
		return nil, apiErr
	}

	// map awss3 objects to file objects
	var objectVersions []*file.ObjectVersion
	for _, v := range response.Versions {
		objectVersions = append(objectVersions, file.NewObjectVersion(v.ETag, v.IsLatest, v.Key, v.LastModified, v.Size, v.VersionId))
	}

	return file.NewStatResponse(objectVersions), nil
}

func (f *fileService) Cat(scope *middleware.Scope, req *file.CatRequest) (*file.CatResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// build request
	request := &awss3.CatRequest{
		Path: req.Path,
		EncryptionKey:       scope.Session.EncryptionKey,
		EncryptionAlgorithm: config.GetS3EncryptionAlgorithm(),
	}

	// call cat
	response, apiErr := scope.AwsS3Provider.Cat(request)
	if apiErr != nil {
		return nil, apiErr
	}

	return file.NewCatResponse(response.File, mime_utils.GetMimeType(req.Path)), nil
}

func (f *fileService) MkDir(scope *middleware.Scope, req *file.MkDirRequest) *error_utils.ApiError {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return reqErr
	}

	// build request
	request := &awss3.MkDirRequest{
		Path: req.Path,
	}

	// call mkDir
	apiErr := scope.AwsS3Provider.MkDir(request)
	if apiErr != nil {
		return apiErr
	}

	return nil
}

func (f *fileService) Rm(scope *middleware.Scope, req *file.RmRequest) *error_utils.ApiError{
	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return reqErr
	}

	// build request
	request := &awss3.RmRequest{
		Path: req.Path,
	}

	// call rm
	apiErr := scope.AwsS3Provider.Rm(request)
	if apiErr != nil {
		return apiErr
	}

	return nil
}

func (f *fileService) Mv(scope *middleware.Scope, req *file.MvRequest) *error_utils.ApiError {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return reqErr
	}

	// build request
	request := &awss3.MvRequest{
		PathFrom:            req.PathFrom,
		PathTo:              req.PathTo,
		EncryptionKey:       scope.Session.EncryptionKey,
		EncryptionAlgorithm: config.GetS3EncryptionAlgorithm(),
	}

	// call move
	apiErr := scope.AwsS3Provider.Mv(request)
	if apiErr != nil {
		return apiErr
	}

	return nil
}

func (f *fileService) Save(scope *middleware.Scope, req *file.SaveRequest) (*file.SaveResponse, *error_utils.ApiError) {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return nil, reqErr
	}

	// convert file header to file
	mFile, err := req.File.Open()
	if err != nil {
		return nil, error_utils.NewInternalServerError(err.Error())
	}

	defer mFile.Close()

	fmt.Printf("Session is %v", scope)

	// build request
	request := &awss3.SaveRequest{
		Path:                req.Path,
		File:                mFile,
		EncryptionKey:       scope.Session.EncryptionKey,
		EncryptionAlgorithm: config.GetS3EncryptionAlgorithm(),
	}

	// call save
	resp, apiErr := scope.AwsS3Provider.Save(request)
	if apiErr != nil {
		return nil, apiErr
	}

	return file.NewSaveResponse(resp.Location, resp.VersionID, resp.UploadID), nil
}

func (f *fileService) Touch(scope *middleware.Scope, req *file.TouchRequest) *error_utils.ApiError {

	// validate request
	reqErr := req.Validate()
	if reqErr != nil {
		return reqErr
	}

	// build request
	request := &awss3.TouchRequest{
		Path:                req.Path,
		EncryptionKey:       scope.Session.EncryptionKey,
		EncryptionAlgorithm: config.GetS3EncryptionAlgorithm(),
	}

	// call touch
	if apiErr := scope.AwsS3Provider.Touch(request); apiErr != nil {
		return apiErr
	}

	return nil
}
