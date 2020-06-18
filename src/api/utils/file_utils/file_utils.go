package file_utils

import (
	"os"
	"time"
)

const (
	DIRECTORY = "directory"
	FILE      = "file"
)

type File struct {
	FName     string `json:"name"`
	FType     string `json:"type"`
	FTime     int64  `json:"time"`
	FSize     int64  `json:"size"`
	FPath     string `json:"path,omitempty"`
	CanRename *bool  `json:"can_rename,omitempty"`
	CanMove   *bool  `json:"can_move_directory,omitempty"`
	CanDelete *bool  `json:"can_delete,omitempty"`
}

func (f File) Name() string {
	return f.FName
}
func (f File) Size() int64 {
	return f.FSize
}
func (f File) Mode() os.FileMode {
	if f.IsDir() {
		return os.ModeDir
	}
	return 0
}
func (f File) ModTime() time.Time {
	if f.FTime == 0 {
		return time.Now()
	}
	return time.Unix(f.FTime, 0)
}
func (f File) IsDir() bool {
	if f.FType != DIRECTORY {
		return false
	}
	return true
}
func (f File) Sys() interface{} {
	return nil
}

type Metadata struct {
	CanSee             *bool      `json:"can_read,omitempty"`
	CanCreateFile      *bool      `json:"can_create_file,omitempty"`
	CanCreateDirectory *bool      `json:"can_create_directory,omitempty"`
	CanRename          *bool      `json:"can_rename,omitempty"`
	CanMove            *bool      `json:"can_move,omitempty"`
	CanUpload          *bool      `json:"can_upload,omitempty"`
	CanDelete          *bool      `json:"can_delete,omitempty"`
	CanShare           *bool      `json:"can_share,omitempty"`
	HideExtension      *bool      `json:"hide_extension,omitempty"`
	RefreshOnCreate    *bool      `json:"refresh_on_create,omitempty"`
	Expire             *time.Time `json:"-"`
}

type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	Time int64  `json:"time"`
}
