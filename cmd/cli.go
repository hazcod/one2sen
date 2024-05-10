package main

import (
	"context"
	"flag"
	"github.com/hazcod/one2sen/config"
	"github.com/hazcod/one2sen/pkg/onepassword"
	msSentinel "github.com/hazcod/one2sen/pkg/sentinel"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "config.yml", "The YAML configuration file.")
	flag.Parse()

	conf := config.Config{}
	if err := conf.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := conf.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.SetLevel(logrusLevel)

	//

	onePass, err := onepassword.New(logger, conf.OnePassword.EventsURL, conf.OnePassword.ApiToken)
	if err != nil {
		logger.WithError(err).Fatal("could not create onepassword client")
	}

	sentinel, err := msSentinel.New(logger, msSentinel.Credentials{
		TenantID:       conf.Microsoft.TenantID,
		ClientID:       conf.Microsoft.AppID,
		ClientSecret:   conf.Microsoft.SecretKey,
		SubscriptionID: conf.Microsoft.SubscriptionID,
		ResourceGroup:  conf.Microsoft.ResourceGroup,
		WorkspaceName:  conf.Microsoft.WorkspaceName,
	})
	if err != nil {
		logger.WithError(err).Fatal("could not create MS Sentinel client")
	}

	//

	if conf.Microsoft.UpdateTable {
		if err := sentinel.CreateTable(ctx, logger, conf.Microsoft.RetentionDays); err != nil {
			logger.WithError(err).Fatal("failed to create MS Sentinel table")
		}
	}

	//

	signinEvents, err := onePass.GetSigninEvents(conf.OnePassword.LookbackDays)
	if err != nil {
		logger.WithError(err).Fatal("could not fetch onepassword signin events")
	}

	signinLogs, err := onepassword.ConvertEventToMap(logger, signinEvents)
	if err != nil {
		logger.WithError(err).Errorf("could not parse signin events")
	}

	//

	usageEvents, err := onePass.GetUsage(conf.OnePassword.LookbackDays)
	if err != nil {
		logger.WithError(err).Fatal("could not fetch onepassword usage events")
	}

	usageLogs, err := onepassword.ConvertUsageToMap(logger, usageEvents)
	if err != nil {
		logger.WithError(err).Error("could not parse usage logs")
	}

	//

	auditEvents, err := onePass.GetAuditEvents(conf.OnePassword.LookbackDays)
	if err != nil {
		logger.WithError(err).Fatal("could not fetch onepassword audit events")
	}

	auditLogs, err := onepassword.ConvertAuditEventToMap(logger, auditEvents)
	if err != nil {
		logger.WithError(err).Errorf("could not parse audit events")
	}

	//

	allLogs := append(signinLogs, usageLogs...)
	allLogs = append(allLogs, auditLogs...)

	logger.WithField("total", len(allLogs)).Info("collected all 1Password logs")

	//

	if err := sentinel.SendLogs(ctx, logger,
		conf.Microsoft.DataCollection.Endpoint,
		conf.Microsoft.DataCollection.RuleID,
		conf.Microsoft.DataCollection.StreamName,
		allLogs); err != nil {
		logger.WithError(err).Fatal("could not ship logs to sentinel")
	}

	//

	logger.WithField("total", len(allLogs)).Info("successfully sent logs to sentinel")
}
