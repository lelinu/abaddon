package file

import (
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/api_utils/utils/validator_utils"
	"time"
)

type StatRequest struct {
	Path string `form:"path" json:"path"`
}

func (req *StatRequest) Validate() *error_utils.ApiError {

	v := validator_utils.NewValidator()

	v.IsNotEmpty("path", req.Path)

	//check if valid
	if !v.IsValid() {
		return error_utils.NewBadRequestError(v.Err.Error())
	}

	return nil
}

type StatResponse struct {
	Versions []*ObjectVersion `json:"versions"`
}

type ObjectVersion struct {
	ETag         *string    `json:"etag"`
	IsLatest     *bool      `json:"is_latest"`
	Key          *string    `json:"key"`
	LastModified *time.Time `json:"last_modified"`
	Size         *int64     `json:"size"`
	VersionId    *string    `json:"version_id"`
}

func NewObjectVersion(etag *string, isLatest *bool, key *string, lastModified *time.Time,
	size *int64, versionId *string) *ObjectVersion{

	return &ObjectVersion{
		ETag:         etag,
		IsLatest:     isLatest,
		Key:          key,
		LastModified: lastModified,
		Size:         size,
		VersionId:    versionId,
	}
}

func NewStatResponse(versions []*ObjectVersion) *StatResponse {
	return &StatResponse{Versions: versions}
}
