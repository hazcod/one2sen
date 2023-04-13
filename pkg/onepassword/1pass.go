package onepassword

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	maxFetch = 100
)

var (
	httpClient = http.Client{Timeout: time.Second * 10}
)

type OnePassword struct {
	Logger   *logrus.Logger
	apiToken string
}

func New(l *logrus.Logger, apiToken string) (*OnePassword, error) {
	if apiToken == "" {
		return nil, errors.New("empty api token provided")
	}

	onePass := OnePassword{
		Logger:   l,
		apiToken: apiToken,
	}

	return &onePass, nil
}
