package onepassword

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"one2sentinel/pkg/utils"
)

const (
	maxFetch = 100
)

type OnePassword struct {
	Logger     *logrus.Logger
	apiToken   string
	httpClient *http.Client
}

func New(l *logrus.Logger, apiToken string) (*OnePassword, error) {
	if apiToken == "" {
		return nil, errors.New("empty api token provided")
	}

	onePass := OnePassword{
		Logger:     l,
		apiToken:   apiToken,
		httpClient: utils.NewLogHttpClient(l),
	}

	return &onePass, nil
}
