package onepassword

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	maxFetch = 100
)

type OnePassword struct {
	Logger     *logrus.Logger
	apiToken   string
	httpClient http.Client
}

type loggingTransport struct {
	transport http.RoundTripper
	Logger    *logrus.Logger
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Logger.IsLevelEnabled(logrus.TraceLevel) {
		// Dump the request and log it
		requestDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			t.Logger.Errorf("Error dumping request: %v", err)
		}
		t.Logger.Tracef("Request:\n%s", string(requestDump))
	}

	// Make the HTTP request using the underlying transport
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if t.Logger.IsLevelEnabled(logrus.TraceLevel) {
		// Dump the response and log it
		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Logger.Errorf("Error dumping response: %v", err)
		}
		t.Logger.Tracef("Response:\n%s", string(responseDump))
	}

	return resp, nil
}

func New(l *logrus.Logger, apiToken string) (*OnePassword, error) {
	if apiToken == "" {
		return nil, errors.New("empty api token provided")
	}

	onePass := OnePassword{
		Logger:     l,
		apiToken:   apiToken,
		httpClient: http.Client{Timeout: time.Second * 10},
	}

	onePass.httpClient.Transport = &loggingTransport{
		Logger:    l,
		transport: http.DefaultTransport,
	}

	return &onePass, nil
}
