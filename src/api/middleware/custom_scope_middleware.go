package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lelinu/abaddon/src/api/config"
	"github.com/lelinu/abaddon/src/api/domain/awss3"
	"github.com/lelinu/abaddon/src/api/providers/awss3_provider"
	"github.com/lelinu/abaddon/src/api/utils/hash_utils"
)

type Scope struct {
	RequestID     string
	Session       Session
	AwsS3Provider awss3_provider.IAwsS3Provider
}

type Session struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	EncryptionKey string `json:"encryption_key"`
	Path          string `json:"path"`
	Region        string `json:"region"`
	Timestamp     string `json:"timestamp"`
}

func NewSession(username string, password string, path string, encryptionKey string,
	region string, timestamp string) *Session {

	return &Session{
		Username:      username,
		Password:      password,
		EncryptionKey: encryptionKey,
		Path:          path,
		Region:        region,
		Timestamp:     timestamp}
}

func newScope(c *gin.Context) *Scope {

	requestID := c.Request.Header.Get("X-Request-Id")
	claims := extractClaims(c)

	// get session token
	session, ok := claims["session"].(string)
	if !ok {
		return nil
	}

	// decrypt session
	decryptedString, err := hash_utils.DecryptString(config.GetApiSecretKey(), session)
	if err != nil {
		fmt.Printf("err is %v", err)
		return nil
	}

	// unmarshal to session
	var sess Session
	err = json.Unmarshal([]byte(decryptedString), &sess)
	if err != nil {
		return nil
	}

	// init request
	initRequest := awss3.NewInitRequest(sess.Username, sess.Password, sess.Path, sess.Region)
	fmt.Printf("Init request is %v", initRequest)

	service, apiErr := awss3_provider.NewAwsS3Provider(config.IsS3LdapEnabled(), config.GetS3EndPoint()).Init(initRequest)
	if apiErr != nil {
		fmt.Printf("error is %v", apiErr)
		return nil
	}

	// return scope
	return &Scope{
		RequestID:     requestID,
		Session:       sess,
		AwsS3Provider: service,
	}
}

// GetScope returns the scope of the current request
func GetScope(c *gin.Context) *Scope {
	return newScope(c)
}
