package awss3

import (
	"encoding/json"
	"github.com/lelinu/api_utils/utils/base64_utils"
)

type RgwToken struct {
	Token *Token `json:"RGW_TOKEN"`
}

type Token struct {
	Version int    `json:"version"`
	Type    string `json:"type"`
	Id      string `json:"id"`
	Key     string `json:"key"`
}

func NewRgwToken(username string, password string) string {
	rgwToken := &RgwToken{
		Token: &Token{
			Version: 1,
			Type:    "ldap",
			Id:      username,
			Key:     password,
		},
	}

	value, _ := json.Marshal(rgwToken)
	return base64_utils.EncodeFromBytes(value)
}
