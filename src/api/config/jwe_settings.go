package config

import (
	"github.com/lelinu/abaddon/src/api/utils/hash_utils"
	"github.com/lelinu/api_utils/utils/env_utils"
	"strconv"
	"time"
)

var (
	keyJweIssuer                    = "JWE_ISSUER"
	keyJweSecretKey                 = "JWE_SECRET_KEY"
	keyJweEncryptionAlgorithm       = "JWE_ENCRYPTION_ALGORITHM"
	keyJweTokenExpiryInHours        = "JWE_TOKEN_EXPIRY_IN_HOURS"
	keyJweRefreshTokenExpiryInHours = "JWE_REFRESH_TOKEN_EXPIRY_IN_HOURS"

	valueJweIssuerFallback                    = "https://golelinu.com"
	valueJweSecretKeyFallback                 = "[NH>!y^Gk-Y36k<sd)bpEZEn$4sLpAbd"
	valueJweEncryptionAlgorithmFallback       = "A256GCM"
	valueJweTokenExpiryInHoursFallback        = "1"
	valueJweRefreshTokenExpiryInHoursFallback = "1"
)

func GetJweIssuer() string {
	return env_utils.GetEnv(keyJweIssuer, valueJweIssuerFallback)
}

func GetJweSecretKey() string {
	secret := env_utils.GetEnv(keyJweSecretKey, valueJweSecretKeyFallback)
	return hash_utils.Hash("JWE_"+secret, len(secret))
}

func GetJweEncryptionAlgorithm() string {
	return env_utils.GetEnv(keyJweEncryptionAlgorithm, valueJweEncryptionAlgorithmFallback)
}

func GetJweTokenExpiryInHours() time.Duration {
	value, _ := strconv.ParseInt(env_utils.GetEnv(keyJweTokenExpiryInHours, valueJweTokenExpiryInHoursFallback), 10, 64)
	return time.Duration(value)
}

func GetJweRefreshTokenExpiryInHours() time.Duration {
	value, _ := strconv.ParseInt(env_utils.GetEnv(keyJweRefreshTokenExpiryInHours, valueJweRefreshTokenExpiryInHoursFallback), 10, 64)
	return time.Duration(value)
}
