package clients

import (
	"github.com/go-redis/redis/v8"
	"github.com/stellar/go/support/log"
	"time"
)

func NewRedis(log *log.Entry) *redis.Client {
	rdsOpts, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		log.WithError(err).Error("failed to parse redis URL")

		rdsOpts = &redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	}
	rdsOpts.ReadTimeout = 500 * time.Millisecond
	rdsOpts.DialTimeout = 500 * time.Millisecond
	rdsOpts.MaxRetries = 1
	rdsOpts.MaxRetryBackoff = 100 * time.Millisecond
	return redis.NewClient(rdsOpts)
}
