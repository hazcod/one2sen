package sentinel

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	insights "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	// has to end with _CL
	tableName = "OnePasswordLogs_CL"
)

func (s *Sentinel) CreateTable(ctx context.Context, l *logrus.Logger, retentionDays uint32) error {
	logger := l.WithField("module", "sentinel_vuln")

	tablesClient, err := insights.NewTablesClient(s.creds.SubscriptionID, s.azCreds, nil)
	if err != nil {
		return fmt.Errorf("could not create ms graph table client: %v", err)
	}

	retention := int32(retentionDays)

	logger.WithField("table_name", tableName).Info("creating or updating table")

	if _, err = tablesClient.Migrate(ctx, s.creds.ResourceGroup, s.creds.WorkspaceName, tableName, nil); err != nil {
		logger.WithError(err).Debug("could not migrate table")
	}

	poller, err := tablesClient.BeginCreateOrUpdate(ctx,
		s.creds.ResourceGroup, s.creds.WorkspaceName, tableName,
		insights.Table{
			Properties: &insights.TableProperties{
				RetentionInDays:      &retention,
				TotalRetentionInDays: to.Ptr[int32](retention * 2),
				Schema: &insights.Schema{
					Columns: []*insights.Column{
						{
							Name: to.Ptr[string]("TimeGenerated"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDateTime),
						},
						{
							Name: to.Ptr[string]("LogType"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("User"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDynamic),
						},
						{
							Name: to.Ptr[string]("Client"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDynamic),
						},
						{
							Name: to.Ptr[string]("Location"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDynamic),
						},
						{
							Name: to.Ptr[string]("Data"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDynamic),
						},
					},
					Name:        to.Ptr[string](tableName),
					Description: to.Ptr[string]("Table that contains events ingested from 1Password."),
				},
			},
		}, nil)
	if err != nil {
		return fmt.Errorf("could not create table '%s': %v", tableName, err)
	}

	_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		return fmt.Errorf("could not poll table creation: %v", err)
	}

	logger.WithField("table_name", tableName).Info("created table")

	return nil
}
