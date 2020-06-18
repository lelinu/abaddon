package awss3

import (
	"encoding/json"
	"github.com/lelinu/api_utils/utils/base64_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRgwToken(t *testing.T) {

	data := NewRgwToken("username", "password")

	assert.NotNil(t, data)

	decodedData, err := base64_utils.DecodeToBytes(data)
	assert.Nil(t, err)
	assert.NotNil(t, decodedData)

	var rgwToken RgwToken
	err = json.Unmarshal(decodedData, &rgwToken)

	assert.Nil(t, err)
	assert.NotNil(t, rgwToken.Token)
	assert.EqualValues(t, "ldap", rgwToken.Token.Type)
	assert.EqualValues(t, 1, rgwToken.Token.Version)
	assert.EqualValues(t, "username", rgwToken.Token.Id)
	assert.EqualValues(t, "password", rgwToken.Token.Key)
}
