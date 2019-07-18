package auth

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	// UserSOCKS5Key defines the Redis key format that will be later formatted into username.
	UserSOCKS5Key = "iap:auth:socks5:%s:password"
	// UserExpiration is time duration before the username and password are expired in Redis.
	UserExpiration = time.Hour * 8

	lowerCase = "abcdefghijklmnopqrstuvwxyz"
	upperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
	specials  = "!#$%'()*+,-./:;=?@[]^_`{|}~"
)

// Client is a struct capable to storing and validating users against temporary data in Redis.
type Client struct {
	store  *redis.Client
	logger *logrus.Logger
}

// New will construct the struct elsewhere.
func New(store *redis.Client, logger *logrus.Logger) *Client {
	return &Client{
		store:  store,
		logger: logger,
	}
}

// Generate will setup a new username and password in the store and return the values.
func (a Client) Generate() (string, string, error) {
	username := randomString(16, upperCase, lowerCase, numbers)
	password := randomString(32, upperCase, lowerCase, numbers, specials)

	err := a.store.Set(fmt.Sprintf(UserSOCKS5Key, username), password, UserExpiration).Err()
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

// Valid will attempt to find the user and compare their passphrase to find a match.
func (a Client) Valid(user, password string) bool {
	a.logger.WithField("user", user).Debugln("authenticating")

	pass, err := a.store.Get(fmt.Sprintf(UserSOCKS5Key, user)).Result()
	if err != nil {
		a.logger.WithField("user", user).Warningln("user not found")
		return false
	}

	return password == pass
}

func randomString(length int, variety ...string) string {
	letter := []rune(strings.Join(variety, ""))

	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
