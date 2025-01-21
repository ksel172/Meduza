package conf

import (
	"fmt"
	"log"
	"os"

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
	MeduzaAdminSecretKeyEnvVar  = "ADMIN_SECRET"
	MeduzaJWTTokenEnvVar        = "JWT_TOKEN"
	BaseConfPathEnvVar          = "BASECONF_PATH"
	BaseConfPathDefault         = "./agent/Agent/baseconf.json"
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

func GetMeduzaAdminSecret() string {
	envToken, ok := os.LookupEnv(MeduzaAdminSecretKeyEnvVar)
	if !ok {
		log.Fatalf("server not configured correctly, missing '%s' environment variable", MeduzaAdminSecretKeyEnvVar)
	}

	return envToken
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
