package internal

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type Context struct {
	Logger *log.Logger
	Redis  *redis.Client
}
