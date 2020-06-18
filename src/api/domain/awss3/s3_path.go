package awss3

import (
	"runtime"
	"strings"
)

var (
	WINDOWS = "windows"
)

type Path struct {
	Bucket string
	Path   string
}

func ParsePath(p string) *Path {

	sep := `/`
	if runtime.GOOS == WINDOWS {
		p = strings.Replace(p, `\`, `/`, -1)
	}

	sp := strings.Split(p, sep)
	bucket := ""
	if len(sp) > 1 {
		bucket = sp[1]
	}
	path := ""
	if len(sp) > 2 {
		path = strings.Join(sp[2:], sep)
	}

	return &Path{
		Bucket: bucket,
		Path:   path,
	}
}
