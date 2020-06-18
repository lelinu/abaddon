package config

import (
	"fmt"
	"github.com/lelinu/abaddon/src/api/utils/hash_utils"
	"github.com/lelinu/abaddon/src/api/utils/path_utils"
	"github.com/lelinu/api_utils/utils/env_utils"
	"os"
	"path/filepath"
	"strconv"
)

var (
	keyApiIsProd      = "API_IS_PROD"
	keyApiVersionUrl  = "API_VERSION_URL"
	keyApiPort        = "API_PORT"
	keyApiSecretKey   = "API_SECRET_KEY"
	keyApiLogPath     = "API_LOG_PATH"
	keyApiLogFileName = "API_LOG_FILE_NAME"

	valueApiIsProdFallback      = "false"
	valueApiVersionUrlFallback  = "api/v1"
	valueApiApiPortFallback     = "8080"
	valueApiSecretKeyFallback   = "-U-ef+d@zfyah~47"
	valueApiLogPathFallback     = "data/log/"
	valueApiLogFileNameFallback = "abaddon_access.log"
)

func init() {
	err := os.MkdirAll(filepath.Join(path_utils.GetCurrentDir(), valueApiLogPathFallback), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func GetApiIsProd() bool {
	value := env_utils.GetEnv(keyApiIsProd, valueApiIsProdFallback)
	convertedValue, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return convertedValue
}

func GetApiVersionUrl() string {
	return env_utils.GetEnv(keyApiVersionUrl, valueApiVersionUrlFallback)
}

func GetApiPort() string {
	return env_utils.GetEnv(keyApiPort, valueApiApiPortFallback)
}

func GetApiSecretKey() string {
	secret := env_utils.GetEnv(keyApiSecretKey, valueApiSecretKeyFallback)
	return hash_utils.Hash("USER_"+secret, len(secret))
}

func GetApiLogPath() string {

	logPath := filepath.Join(path_utils.GetCurrentDir(), valueApiLogPathFallback)
	fullLogFilePath := filepath.Join(logPath, valueApiLogFileNameFallback)

	_, err := os.OpenFile(fullLogFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Log file path is %v \n", fullLogFilePath)

	return fullLogFilePath
}
