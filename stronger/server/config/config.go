package config

import (
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/config"
)

type Server struct {
	Postgresql  config.Postgresql `mapstructure:"postgresql" json:"postgresql" yaml:"postgresql"`
	UserDB      config.Postgresql `mapstructure:"user-db" json:"userDB" yaml:"user-db"`
	Casbin      Casbin            `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	Redis       Redis             `mapstructure:"redis" json:"redis" yaml:"redis"`
	System      System            `mapstructure:"system" json:"system" yaml:"system"`
	JWT         JWT               `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Captcha     Captcha           `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Zap         Zap               `mapstructure:"zap" json:"zap" yaml:"zap"`
	LocalUpload LocalUpload       `mapstructure:"localUpload" json:"localUpload" yaml:"localUpload"`
	Couchdb     Couchdb           `mapstructure:"couchdb" json:"couchdb" yaml:"couchdb"`
	MongoDB     MongoDB           `mapstructure:"mongoDB" json:"mongoDB" yaml:"mongoDB"`
	//FilePath    FilePath    `mapstructure:"filepath" json:"filepath" yaml:"filepath"`
}

//type FilePath struct {
//	ConstXml string `mapstructure:"const-xml" json:"constXml" yaml:"const-xml"`
//}

type Couchdb struct {
	Url       string `mapstructure:"url" json:"url" yaml:"url"`
	Username  string `mapstructure:"username" json:"username" yaml:"username"`
	Password  string `mapstructure:"password" json:"password" yaml:"password"`
	ExpectURL string `mapstructure:"expect-url" json:"expectUrl" yaml:"expect-url"`
}

type System struct {
	UseMultipoint   bool     `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	Env             string   `mapstructure:"env" json:"env" yaml:"env"`
	Port            int      `mapstructure:"port" json:"port" yaml:"port"`
	CommonPort      int      `mapstructure:"common-port" json:"commonPort" yaml:"common-port"`
	DbType          string   `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	Name            string   `mapstructure:"name" json:"name" yaml:"name"`
	ProCode         string   `mapstructure:"pro-code" json:"proCode" yaml:"pro-code"`
	Process         string   `mapstructure:"process" json:"process" yaml:"process"`
	ImageEncrypt    bool     `mapstructure:"image-encrypt" json:"imageEncrypt" yaml:"image-encrypt"`
	DownloadPath    string   `mapstructure:"downloadPath" json:"downloadPath" yaml:"downloadPath"`
	DownloadBillEnd string   `mapstructure:"downloadBillEnd" json:"downloadBillEnd" yaml:"downloadBillEnd"`
	ProArr          []string `mapstructure:"pro-arr" json:"proArr" yaml:"pro-arr"`
	ConstUrl        string   `mapstructure:"const-url" json:"constUrl" yaml:"const-url"`
}

type JWT struct {
	SigningKey string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`
	ExpiresAt  int64  `mapstructure:"expires-at" json:"expiresAt" yaml:"expires-at"`
}

type Casbin struct {
	ModelPath string `mapstructure:"model-path" json:"modelPath" yaml:"model-path"`
}

type Postgresql struct {
	Username             string `mapstructure:"username" json:"username" yaml:"username"`
	Password             string `mapstructure:"password" json:"password" yaml:"password"`
	Host                 string `mapstructure:"host" json:"host" yaml:"host"`
	Dbname               string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Port                 string `mapstructure:"port" json:"port" yaml:"port"`
	Config               string `mapstructure:"config" json:"config" yaml:"config"`
	MaxIdleConns         int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns         int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	Logger               bool   `mapstructure:"logger" json:"logger" yaml:"logger"`
	PreferSimpleProtocol bool   `mapstructure:"prefer-simple-protocol" json:"preferSimpleProtocol" yaml:"prefer-simple-protocol"`
}

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

type LocalUpload struct {
	Local    bool   `mapstructure:"local" json:"local" yaml:"local"`
	FilePath string `mapstructure:"file-path" json:"filePath" yaml:"file-path"`
}
type Captcha struct {
	KeyLong       int           `mapstructure:"key-long" json:"keyLong" yaml:"key-long"`
	ImgWidth      int           `mapstructure:"img-width" json:"imgWidth" yaml:"img-width"`
	ImgHeight     int           `mapstructure:"img-height" json:"imgHeight" yaml:"img-height"`
	Expiration    time.Duration `mapstructure:"expiration" json:"expiration" yaml:"expiration"`
	GCLimitNumber int           `mapstructure:"gc-limit-number" json:"GCLimitNumber" yaml:"gc-limit-number"`
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	LinkName      string `mapstructure:"link-name" json:"linkName" yaml:"link-name"`
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"showLine"`
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`
}

type MongoDB struct {
	Url string `json:"url" yaml:"url"`
}
