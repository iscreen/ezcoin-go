package config

type System struct {
	Env             string `mapstructure:"env" json:"env" yaml:"env"`
	Addr            int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType          string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	OssType         string `mapstructure:"oss-type" json:"ossType" yaml:"oss-type"`
	UseMultipoint   bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	LimitCountIP    int    `mapstructure:"iplimit-count" json:"iplimitCount" yaml:"iplimit-count"`
	LimitTimeIP     int    `mapstructure:"iplimit-time" json:"iplimitTime" yaml:"iplimit-time"`
	DefaultPageSize int    `mapstructure:"default_page_size" json:"defaultPageSize" yaml:"default_page_size"`
	MaxPageSize     int    `mapstructure:"max_page_size" json:"maxPageSize" yaml:"max_page_size"`
	FrontUrl        string `mapstructure:"front-url" json:"frontUrl" yaml:"front-url"`
}
