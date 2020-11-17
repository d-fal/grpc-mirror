package identity

import (
	"fmt"
	"regexp"
)

var (
	signatureObject Parameters
	// ConfigList list of errors

	r    *regexp.Regexp
	trim *regexp.Regexp
)

// Parameters gets app parameters
type Parameters struct {
	App App `yaml:"App"`
}

// App app struct
type App struct {
	Name    string     `yaml:"Name"`
	ID      string     `yaml:"ID"`
	Version float64    `yaml:"Version"`
	Routes  []Receptor `yaml:"Routes"`
	BaseURL string     `yaml:"BaseURL"`
}

// Receptor receptor that handles service
type Receptor struct {
	ID          string       `yaml:"HandlerID"`
	Name        string       `yaml:"Name"`
	Method      string       `yaml:"Method"`
	URL         string       `yaml:"URL"`
	Middlewares []Middleware `yaml:"Middleware"`
}

// Middleware that specifies the middlewares an app uses
type Middleware struct {
	ID  string `yaml:"ID"`
	URL string `yaml:"URL"`
}

func SetSignature(signature Parameters) {

	signatureObject = signature
}

func GetSignature() Parameters {
	return signatureObject
}

func GetAppIssuer() string {
	return fmt.Sprintf("%s, version %2.1f", signatureObject.App.Name, signatureObject.App.Version)
}
