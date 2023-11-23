package sentinel

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	logsPerRequest = 10
)

func chunkLogs(slice []map[string]string, chunkSize int) [][]map[string]string {
	var chunks [][]map[string]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// avoid slicing beyond slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func (s *Sentinel) SendLogs(ctx context.Context, l *logrus.Logger, endpoint, ruleID, streamName string, logs []map[string]string) error {
	logger := l.WithField("module", "sentinel_logs")

	logger.WithField("table_name", tableName).WithField("total", len(logs)).Info("shipping logs")

	chunkedLogs := chunkLogs(logs, logsPerRequest)
	for i, logsChunk := range chunkedLogs {
		l.WithField("progress", fmt.Sprintf("%d/%d", i+1, len(chunkedLogs))).Debug("ingesting log chunks")

		if len(logsChunk) == 0 {
			l.Warn("processing empty chunk")
			continue
		}

		if err := s.IngestLog(ctx, endpoint, ruleID, streamName, logsChunk); err != nil {
			return fmt.Errorf("could not ingest log: %v", err)
		}
	}

	//

	logger.WithField("table_name", tableName).Info("shipped logs")

	return nil
}
