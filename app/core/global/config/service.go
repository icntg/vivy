package config

import "fmt"

type Service struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         uint16 `mapstructure:"port" json:"port" yaml:"port"`
	StaticPath   string `mapstructure:"static-path" json:"staticPath" yaml:"static-path"`
	TemplateHtml bool   `mapstructure:"template-html" json:"templateHtml" yaml:"template-html"`
}

func GetServiceAddress(ths *Service) string {
	return fmt.Sprintf("%s:%d", ths.Host, ths.Port)
}
