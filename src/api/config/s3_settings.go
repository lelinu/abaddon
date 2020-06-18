package config

import (
	"github.com/lelinu/api_utils/utils/env_utils"
	"strconv"
)

var (
	//keys
	keyS3Region       = "S3_REGION"
	keyS3EncAlgorithm = "S3_ENC_ALGORITHM"
	keyS3LdapEnabled  = "S3_LDAP_ENABLED"
	keyS3EndPoint     = "S3_ENDPOINT"

	//fall backs
	valueS3RegionFallback       = "eu-west-1"
	valueS3EncAlgorithmFallback = "AES256"
	valueS3LdapEnabled          = "true"
	valueS3EndPoint             = "https://api.storage.lelinu.com"
)

func GetS3Region() string {
	return env_utils.GetEnv(keyS3Region, valueS3RegionFallback)
}

func GetS3EncryptionAlgorithm() string {
	return env_utils.GetEnv(keyS3EncAlgorithm, valueS3EncAlgorithmFallback)
}

func GetS3EndPoint() string {
	return env_utils.GetEnv(keyS3EndPoint, valueS3EndPoint)
}

func IsS3LdapEnabled() bool {
	value := env_utils.GetEnv(keyS3LdapEnabled, valueS3LdapEnabled)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return boolValue
}

