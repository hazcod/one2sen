package sentinel

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	ingestTimeout = time.Second * 10
)

func (s *Sentinel) IngestLog(ctx context.Context, logs []map[string]string) error {
	/*
		ingest, err := azingest.NewClient("1password", s.azCreds, nil)
		if err != nil {
			return fmt.Errorf("could not create azure ingest client: %v", err)
		}

		ingest.Upload(ctx)
	*/

	logger := s.logger.WithField("module", "sentinel_ingest")

	logPayload, err := json.Marshal(&logs)
	if err != nil {
		return fmt.Errorf("could not json encode log message: %v", err)
	}

	if s.logger.IsLevelEnabled(logrus.TraceLevel) {
		fmt.Println(string(logPayload))
	}

	// prep timestamp
	dateString := time.Now().UTC().Format(time.RFC1123)
	dateString = strings.Replace(dateString, "UTC", "GMT", -1)

	// build log request signature
	hashString := fmt.Sprintf("POST\n%d\napplication/json\nx-ms-date:%s\n/api/logs", len(logPayload), dateString)
	hashedString, err := BuildSignature(hashString, s.creds.WorkspaceKey)
	if err != nil {
		return fmt.Errorf("could not build log signature: %v", err)
	}

	signature := fmt.Sprintf("SharedKey %s:%s", s.creds.WorkspaceID, hashedString)
	url := fmt.Sprintf("https://%s.ods.opinsights.azure.com/api/logs?api-version=2016-04-01", s.creds.WorkspaceID)

	ingestCtx, _ := context.WithTimeout(ctx, ingestTimeout)

	req, err := http.NewRequestWithContext(ingestCtx, "POST", url, bytes.NewReader(logPayload))
	if err != nil {
		return fmt.Errorf("could not create http request: %v", err)
	}

	req.Header.Add("Log-Type", tableName)
	req.Header.Add("Authorization", signature)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-ms-date", dateString)
	req.Header.Add("time-generated-field", "TimeGenerated")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send log: %v", err)
	}
	defer resp.Body.Close()

	bv, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read ingest response: %v", err)
	}

	logger.WithField("response", string(bv)).WithField("resp_code", resp.StatusCode).Debug("got ingest response")

	if resp.StatusCode != 200 {
		return fmt.Errorf("got status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	logger.WithField("response_code", resp.StatusCode).Debug("successfully shipped log")

	return nil
}

func BuildSignature(message, secret string) (string, error) {
	if message == "" {
		return "", errors.New("empty message")
	}

	if secret == "" {
		return "", errors.New("empty secret")
	}

	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("invalid base64 secret: %v", err)
	}

	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
