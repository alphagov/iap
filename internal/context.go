package internal

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// Context defines what is required in different commands of the application at the same time.
type Context struct {
	Logger *log.Logger
	Redis  *redis.Client
}
