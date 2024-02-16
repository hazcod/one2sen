package config

import "testing"

func TestConfig_Validate(t *testing.T) {
	conf := Config{}
	if err := conf.Validate(); err == nil {
		t.Fail()
	}
}
