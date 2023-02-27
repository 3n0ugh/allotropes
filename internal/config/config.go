package config

import "os"

const (
	secret            = "test"
	postgresqlDSN     = "localdsn"
	couchbaseDSN      = "localdsn"
	couchbaseUsername = "localusername"
	couchbasePassword = "cbpass"
	couchbaseBucket   = "cbbucket"
)

type Config struct {
	Application Application
	PostgreSQL  PostgreSQL
	Couchbase   Couchbase
}

type Application struct {
	Secret string
}

type PostgreSQL struct {
	DataSource string
}

type Couchbase struct {
	BucketName string
	UserName   string
	Password   string
	DataSource string
}

func ReadConfig() Config {
	return Config{
		Application: Application{
			Secret: setConfig("SECRET", secret),
		},
		PostgreSQL: PostgreSQL{
			DataSource: setConfig("POSTGRES_DSN", postgresqlDSN),
		},
		Couchbase: Couchbase{
			BucketName: setConfig("CB_BUCKET", couchbaseBucket),
			UserName:   setConfig("CB_USERNAME", couchbaseUsername),
			Password:   setConfig("CB_PASSWORD", couchbasePassword),
			DataSource: setConfig("CB_DSN", couchbaseDSN),
		},
	}
}

func setConfig(configName, defaultVal string) string {
	if cfg := os.Getenv(configName); cfg != "" {
		return cfg
	}
	return defaultVal
}
