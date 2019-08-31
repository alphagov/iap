package auth

import (
	"crypto/rand"
	"fmt"
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
	username, err := randomString(16, upperCase, lowerCase, numbers)
	if err != nil {
		return "", "", err
	}
	password, err := randomString(32, upperCase, lowerCase, numbers, specials)
	if err != nil {
		return "", "", err
	}

	if err := a.store.Set(fmt.Sprintf(UserSOCKS5Key, username), password, UserExpiration).Err(); err != nil {
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

func randomString(length int, charSet ...string) (string, error) {
	letters := strings.Join(charSet, "")
	bytes, err := generateRandomBytes(length)
	if err != nil {
		return "", fmt.Errorf("unable to generate random string: %s", err)
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func generateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("unable to generate random bytes: %s", err)
	}

	return b, nil
}
