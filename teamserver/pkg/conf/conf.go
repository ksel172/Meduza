package conf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ksel172/Meduza/teamserver/utils"
)

const (
	MeduzaServerHostnameEnvVar  = "TEAMSERVER_HOSTNAME"
	MeduzaServerHostnameDefault = "localhost"
	MeduzaServerPortEnvVar      = "TEAMSERVER_PORT"
	MeduzaServerPortDefault     = 8080
	MeduzaServerModeEnvVar      = "TEAMSERVER_MODE"
	MeduzaServerModeDefault     = "dev"
	MeduzaDbHostnameEnvVar      = "DB_HOST"
	MeduzaDbHostnameDefault     = "localhost"
	MeduzaDbPortEnvVar          = "DB_PORT"
	MeduzaDbPortDefault         = 5432
	MeduzaDbUsernameEnvVar      = "DB_USER"
	MeduzaDbUsernameDefault     = "postgres"
	MeduzaDbPasswordEnvVar      = "DB_PASSWORD"
	MeduzaDbPasswordDefault     = "postgres"
	MeduzaDbNameEnvVar          = "DB_NAME"
	MeduzaDbNameDefault         = "meduza_db"
	MeduzaDbSchemaEnvVar        = "DB_SCHEMA"
	MeduzaDbSchemaDefault       = "meduza_schema"
	MeduzaRedisHostEnvVar       = "REDIS_HOST"
	MeduzaRedisHostDefault      = "localhost"
	MeduzaRedisPortEnvVar       = "REDIS_PORT"
	MeduzaRedisPortDefault      = 6379
	MeduzaRedisPasswordEnvVar   = "REDIS_PASSWORD"
	MeduzaRedisPasswordDefault  = "password"
	MeduzaJWTTokenEnvVar        = "JWT_TOKEN"
	BaseConfPathEnvVar          = "BASECONF_PATH"
	BaseConfPathDefault         = "./agent/Agent/baseconf.json"
	ModuleUploadPathEnvVar      = "MODULE_UPLOAD_PATH"
	ModuleUploadPathDefault     = "./teamserver/modules"
	ListenPortRangeStartEnvVar  = "LISTENER_PORT_RANGE_START"
	ListenPortRangeEndEnvVar    = "LISTENER_PORT_RANGE_END"
	MeduzaCertUploadPathEnvVar  = "CERT_UPLOAD_PATH"
	MeduzaCertUploadPathDefault = "./teamserver/certs"
)

func GetMeduzaServerHostname() string {

	return utils.GetEnvString(MeduzaServerHostnameEnvVar, MeduzaServerHostnameDefault)
}

func GetMeduzaServerPort() int {
	return utils.GetEnvInt(MeduzaServerPortEnvVar, MeduzaServerPortDefault)
}

func GetMeduzaDbHostname() string {
	return utils.GetEnvString(MeduzaDbHostnameEnvVar, MeduzaDbHostnameDefault)
}

func GetMeduzaDbPort() int {
	return utils.GetEnvInt(MeduzaDbPortEnvVar, MeduzaDbPortDefault)
}

func GetMeduzaDbUsername() string {
	return utils.GetEnvString(MeduzaDbUsernameEnvVar, MeduzaDbUsernameDefault)
}

func GetMeduzaDbPassword() string {
	return utils.GetEnvString(MeduzaDbPasswordEnvVar, MeduzaDbPasswordDefault)
}

func GetMeduzaDbName() string {
	return utils.GetEnvString(MeduzaDbNameEnvVar, MeduzaDbNameDefault)
}

func GetMeduzaDbSchema() string {
	return utils.GetEnvString(MeduzaDbSchemaEnvVar, MeduzaDbSchemaDefault)
}

func GetMeduzaRedisAddress() string {
	return fmt.Sprintf("%s:%d", utils.GetEnvString(MeduzaRedisHostEnvVar, MeduzaRedisHostDefault), utils.GetEnvInt(MeduzaRedisPortEnvVar, MeduzaRedisPortDefault))
}

func GetMeduzaRedisPassword() string {
	return utils.GetEnvString(MeduzaRedisPasswordEnvVar, MeduzaRedisPasswordDefault)
}

func GetMeduzaServerMode() string {
	return utils.GetEnvString(MeduzaServerModeEnvVar, MeduzaServerModeDefault)
}

func GetMeduzaJWTToken() string {
	jwt, ok := os.LookupEnv(MeduzaJWTTokenEnvVar)
	if !ok {
		log.Fatalf("server not configured correctly, missing %s environment variable", MeduzaJWTTokenEnvVar)
	}

	return jwt
}

func GetBaseConfPath() string {
	return utils.GetEnvString(BaseConfPathEnvVar, BaseConfPathDefault)
}

func GetModuleUploadPath() string {
	return utils.GetEnvString(ModuleUploadPathEnvVar, ModuleUploadPathDefault)
}

func GetCertUploadPath() string {
	return utils.GetEnvString(MeduzaCertUploadPathEnvVar, MeduzaCertUploadPathDefault)
}

func GetProjectRootPath() string {
	// Get the directory of the currently running file
	execPath, err := os.Executable()
	if err != nil {
		panic("Failed to determine the executable path: " + err.Error())
	}

	// Resolve the executable path to the root project directory
	projectRoot := filepath.Dir(filepath.Dir(execPath)) // Assuming the binary is two levels deep in the project
	return projectRoot
}

func GetListenerPortRangeStart() int {
	return utils.GetEnvInt(ListenPortRangeStartEnvVar, 8000)
}

func GetListenerPortRangeEnd() int {
	return utils.GetEnvInt(ListenPortRangeEndEnvVar, 8010)
}
