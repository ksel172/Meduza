package conf

import (
	"log"
	"os"
	"strconv"
)

const (
	MeduzaServerHostnameEnvVar  = "TEAMSERVER_HOSTNAME"
	MeduzaServerHostnameDefault = "localhost"
	MeduzaServerPortEnvVar      = "TEAMSERVER_PORT"
	MeduzaServerPortDefault     = "8080"
	MeduzaDbHostnameEnvVar      = "DB_HOST"
	MeduzaDbHostnameDefault     = "localhost"
	MeduzaDbPortEnvVar          = "DB_PORT"
	MeduzaDbPortDefault         = "5432"
	MeduzaDbUsernameEnvVar      = "DB_USER"
	MeduzaDbUsernameDefault     = "postgres"
	MeduzaDbPasswordEnvVar      = "DB_PASS"
	MeduzaDbPasswordDefault     = "postgres"
	MeduzaDbNameEnvVar          = "DB_NAME"
	MeduzaDbNameDefault         = "meduza_db"
	MeduzaDbSchemaEnvVar        = "DB_SCHEMA"
	MeduzaDbSchemaDefault       = "meduza_schema"
)

func GetMeduzaServerHostname() string {
	hostname, exists := os.LookupEnv(MeduzaServerHostnameEnvVar)
	if !exists {
		log.Printf("Environment variable '%s' not set, defaulting to '%s'...\n", MeduzaServerHostnameEnvVar, MeduzaServerHostnameDefault)
		hostname = MeduzaServerHostnameDefault
	}
	return hostname
}

func GetMeduzaServerPort() int {
	port, exists := os.LookupEnv(MeduzaServerPortEnvVar)

	if !exists {
		log.Printf("Environment variable '%s' not set, defaulting to '%s'...\n", MeduzaServerPortEnvVar, MeduzaServerPortDefault)
		port = MeduzaServerPortDefault
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("Environmental variable '%s' is not set to a number (%s). Defaulting to '%s'...", MeduzaServerPortEnvVar, port, MeduzaServerPortDefault)
		port = MeduzaServerPortDefault
	}

	return portNumber
}

func GetMeduzaDbHostname() string {
	hostname, exists := os.LookupEnv(MeduzaDbHostnameEnvVar)
	if !exists {
		log.Printf("Environmental variable '%s' is not set, defaulting to '%s'...", MeduzaDbHostnameEnvVar, MeduzaDbHostnameDefault)
		hostname = MeduzaDbHostnameDefault
	}
	return hostname
}

func GetMeduzaDbPort() string {
	port, exists := os.LookupEnv(MeduzaDbPortEnvVar)
	if !exists {
		log.Printf("Environmental variable '%s' is not set, defaulting to '%s'...", MeduzaDbPortEnvVar, MeduzaDbPortDefault)
		port = MeduzaDbPortDefault
	}
	return port
}

func GetMeduzaDbUsername() string {
	username, exists := os.LookupEnv(MeduzaDbUsernameEnvVar)
	if !exists {
		log.Printf("Environmental variable '%s' is not set, defaulting to '%s'...", MeduzaDbUsernameEnvVar, MeduzaDbUsernameDefault)
		username = MeduzaDbUsernameDefault
	}
	return username
}

func GetMeduzaDbPassword() string {
	password, exists := os.LookupEnv(MeduzaDbPasswordEnvVar)
	if !exists {
		log.Printf("Environmental variable '%s' is not set, defaulting to '%s'...", MeduzaDbPasswordEnvVar, MeduzaDbPasswordDefault)
		password = MeduzaDbPasswordDefault
	}
	return password
}

func GetMeduzaDbName() string {
	name, exists := os.LookupEnv(MeduzaDbNameEnvVar)
	if !exists {
		log.Printf("Environmental variable '%s' is not set, defaulting to '%s'...", MeduzaDbNameEnvVar, MeduzaDbNameDefault)
		name = MeduzaDbNameDefault
	}
	return name
}

func GetMeduzaDbSchema() string {
	schema, exists := os.LookupEnv(MeduzaDbSchemaEnvVar)
	if !exists {
		log.Printf("Environmental variable '%s' is not set, defaulting to '%s'...", MeduzaDbSchemaEnvVar, MeduzaDbSchemaDefault)
		schema = MeduzaDbSchemaDefault
	}
	return schema
}
