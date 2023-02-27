package database

import (
	"time"

	"github.com/3n0ugh/allotropes/internal/config"
	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/couchbase/gocb/v2"
)

func OpenConnectionCB(cfg config.Config) (*gocb.Bucket, error) {
	cluster, err := gocb.Connect(cfg.Couchbase.DataSource, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: cfg.Couchbase.UserName,
			Password: cfg.Couchbase.Password,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "couchbase connection")
	}

	bucket := cluster.Bucket(cfg.Couchbase.BucketName)
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		return nil, errors.Wrap(err, "bucket connection")
	}

	return bucket, nil
}
