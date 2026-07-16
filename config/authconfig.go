package config

import "time"

type AuthConfig struct {
	AccessSecret  string        `mapstructure:"access_secret" validate:"required,min=32,excludes=your_super_secret"`
	RefreshSecret string        `mapstructure:"refresh_secret" validate:"required,min=32,excludes=your_super_secret"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}
