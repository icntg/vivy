package config

import "fmt"

type Service struct {
	Host               string `mapstructure:"host" json:"host" yaml:"host"`
	Port               uint16 `mapstructure:"port" json:"port" yaml:"port"`
	StaticPath         string `mapstructure:"static-path" json:"staticPath" yaml:"static-path"`
	TemplateHtml       bool   `mapstructure:"template-html" json:"templateHtml" yaml:"template-html"`
	SessionSecret      string `mapstructure:"session-secret" json:"sessionSecret" yaml:"session-secret"`
	SessionSecretBytes []byte `mapstructure:"-" json:"-" yaml:"-"`
	CookieTimeout      int    `mapstructure:"cookie-timeout" json:"cookieTimeout" yaml:"cookie-timeout"`
}

func (ths *Service) GetServiceAddress() string {
	return fmt.Sprintf("%s:%d", ths.Host, ths.Port)
}
