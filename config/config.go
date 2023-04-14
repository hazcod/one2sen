package config

import (
	"errors"
	"fmt"
	validator "github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultLogLevel      = "DEBUG"
	defaultExpiresMonths = 6
)

type Config struct {
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	} `json:"log"`

	OnePassword struct {
		ApiToken string `yaml:"api_token"`
	} `json:"onepassword"`

	Microsoft struct {
		AppID          string `yaml:"app_id" env:"MS_APP_ID" valid:"minstringlength(3)"`
		SecretKey      string `yaml:"secret_key" env:"MS_SECRET_KEY" valid:"minstringlength(3)"`
		TenantID       string `yaml:"tenant_id" env:"MS_TENANT_ID" valid:"minstringlength(3)"`
		SubscriptionID string `yaml:"subscription_id" env:"MS_SUB_ID" valid:"minstringlength(3)"`
		ResourceGroup  string `yaml:"resource_group" env:"MS_RSG_ID" valid:"minstringlength(3)"`
		WorkspaceName  string `yaml:"workspace_name" env:"MS_WS_NAME" valid:"minstringlength(3)"`
		WorkspaceID    string `yaml:"workspace_id" env:"MS_WS_ID" valid:"minstringlength(3)"`
		WorkspaceKey   string `yaml:"workspace_primary_key" env:"MS_WS_KEY" valid:"minstringlength(10)"`
		ExpiresMonths  uint16 `yaml:"expires_months" env:"MS_EXPIRES_MONTHS"`

		UpdateTable bool `yaml:"update_table" env:"MS_UPDATE_TABLE"`
	} `yaml:"microsoft"`
}

func (c *Config) Validate() error {
	if c.Log.Level == "" {
		c.Log.Level = defaultLogLevel
	}

	if c.OnePassword.ApiToken == "" {
		return errors.New("no onepassword api token provided")
	}

	if c.Microsoft.ExpiresMonths == 0 {
		c.Microsoft.ExpiresMonths = defaultExpiresMonths
	}

	if valid, err := validator.ValidateStruct(c); !valid || err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	return nil
}

func (c *Config) Load(path string) error {
	if path != "" {
		configBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load configuration file at '%s': %v", path, err)
		}

		if err = yaml.Unmarshal(configBytes, c); err != nil {
			return fmt.Errorf("failed to parse configuration: %v", err)
		}
	}

	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("could not load environment: %v", err)
	}

	return nil
}
