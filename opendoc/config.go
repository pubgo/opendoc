package opendoc

import "github.com/pubgo/funk/version"

// https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/

type Config struct {
	Title              string                 `yaml:"title"`
	OpenapiRouter      string                 `yaml:"path"`
	OpenapiRedocRouter string                 `yaml:"redoc-path"`
	OpenapiUrl         string                 `yaml:"openapi-path"`
	OpenapiOpt         map[string]interface{} `yaml:"options"`
}

func DefaultCfg() *Config {
	return &Config{
		Title:              version.Project() + " openapi docs",
		OpenapiRouter:      "/debug/docs",
		OpenapiRedocRouter: "/debug/redocs",
		OpenapiUrl:         "/debug/docs/openapi.yaml",
		OpenapiOpt:         make(map[string]interface{}),
	}
}
