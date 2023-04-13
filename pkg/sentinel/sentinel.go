package sentinel

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type Credentials struct {
	TenantID       string
	ClientID       string
	ClientSecret   string
	SubscriptionID string
	ResourceGroup  string
	WorkspaceName  string
	WorkspaceID    string
	WorkspaceKey   string
}

type Sentinel struct {
	creds Credentials

	azCreds *azidentity.ClientSecretCredential
}

func New(creds Credentials) (*Sentinel, error) {
	sentinel := Sentinel{
		creds: creds,
	}

	if creds.WorkspaceID == "" {
		return nil, errors.New("no workspace id provided")
	}

	azCreds, err := azidentity.NewClientSecretCredential(creds.TenantID, creds.ClientID, creds.ClientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("could not authenticate to MS Sentinel: %v", err)
	}

	sentinel.azCreds = azCreds

	return &sentinel, nil
}
