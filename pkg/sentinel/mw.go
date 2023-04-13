package sentinel

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

func getLoggingHttpClient(l *logrus.Logger) *http.Client {
	httpClient := http.Client{
		Timeout: ingestTimeout,
	}

	httpClient.Transport = taggedRoundTripper{Logger: l}

	return &httpClient
}

type taggedRoundTripper struct {
	Logger *logrus.Logger
}

// RoundTrip injects a http request header on every request and logs request/response
func (t taggedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)

	if t.Logger != nil && t.Logger.IsLevelEnabled(logrus.TraceLevel) {
		dumped, err := httputil.DumpRequest(req, true)
		if err != nil {
			t.Logger.WithError(err).Error("could not dump http request")
		} else {
			t.Logger.Trace(string(dumped))
		}

		if req.Response != nil {
			dumped, err = httputil.DumpResponse(req.Response, true)
			if err != nil {
				t.Logger.WithError(err).Error("could not dump http response")
			} else {
				t.Logger.Trace(string(dumped))
			}
		}
	}

	return resp, err
}
