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
	apiURL     string
}

func New(l *logrus.Logger, tenantURL string, apiToken string) (*OnePassword, error) {
	if apiToken == "" {
		return nil, errors.New("empty api token provided")
	}

	if tenantURL == "" {
		return nil, errors.New("no tenant TLD provided such as com,ca,eu")
	}

	onePass := OnePassword{
		Logger:     l,
		apiToken:   apiToken,
		httpClient: utils.NewLogHttpClient(l),
		apiURL:     tenantURL,
	}

	return &onePass, nil
}
