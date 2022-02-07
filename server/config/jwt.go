package config

type JWT struct {
	SigningKey  string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`
	ExpiresTime int64  `mapstructure:"expires-time" json:"expiresTime" yaml:"expires-time"`
	BufferTime  int64  `mapstructure:"buffer-time" json:"bufferTime" yaml:"buffer-time"`
	Issuer      string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
	TokenName   string `mapstructure:"token-name" json:"tokenName" yaml:"token-name"`
}
