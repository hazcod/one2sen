package sentinel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest"
	"github.com/sirupsen/logrus"
)

func (s *Sentinel) IngestLog(ctx context.Context, endpoint, ruleID, streamName string, logs []map[string]string) error {
	logger := s.logger.WithField("module", "sentinel_ingest")

	ingest, err := azingest.NewClient(endpoint, s.azCreds, nil)
	if err != nil {
		return fmt.Errorf("could not create azure ingest client: %v", err)
	}

	for i, log := range logs {
		logPayload, err := json.Marshal(&log)
		if err != nil {
			return fmt.Errorf("could not json encode log message: %v", err)
		}

		if s.logger.IsLevelEnabled(logrus.TraceLevel) {
			logger.Tracef("%s", string(logPayload))
		}

		logger.WithField("logs", fmt.Sprintf("%d/%d", i+1, len(logs))).Debug("uploading log")

		_, err = ingest.Upload(ctx, ruleID, streamName, logPayload, nil)
		if err != nil {
			return fmt.Errorf("could not upload logs: %v", err)
		}
	}

	logger.WithField("total_logs", len(logs)).Debug("successfully uploaded 1password logs")

	return nil
}
