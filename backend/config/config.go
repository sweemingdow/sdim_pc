package config

import "time"

type Config struct {
	AppCfg     AppConfig     `yaml:"app-config"`
	LogCfg     LogConfig     `yaml:"log-config"`
	HttpReqCfg HttpReqConfig `yaml:"http-req-config"`
}

type AppConfig struct {
	Profile          string        `yaml:"profile"`
	StopTimeout      time.Duration `yaml:"stop-timeout"`
	EngineServerAddr string        `yaml:"engine-server-addr"`
	ApiGwAddr        string        `yaml:"api-gw-addr"`
}

type LogConfig struct {
	Level       string `yaml:"level"`
	ColorfulStd bool   `yaml:"colorful-std"`
	FilePath    string `yaml:"file-path"`
	MaxFileSize int    `yaml:"max-file-size"`
	MaxBackup   int    `yaml:"max-backup"`
	HistoryDays int    `yaml:"history-days"`
	Compress    bool   `yaml:"compress"`
}

type HttpReqConfig struct {
	MaxConn         int           `yaml:"max-conn"`
	MaxIdleConn     int           `yaml:"max-idle-conn"`
	ConnIdleTimeout time.Duration `yaml:"conn-idle-timeout"`
	RequestTimeout  time.Duration `yaml:"request-timeout"`
	ResponseTimeout time.Duration `yaml:"response-timeout"`
}
