package routing

import (
	"encoding/json"
	commonHttp "github.com/kulycloud/common/http"
)

type RouteList []Route

type Config struct {
	Routes RouteList `json:"routes"`
}

type Route struct {
	Path        *PathSpecification `json:"path,omitempty"`
	Methods     *[]string          `json:"methods,omitempty"`
	Destination *string            `json:"destination,omitempty"`
	Subroutes   *RouteList           `json:"subroutes"`
}

type PathSpecification struct {
	Prefix  *string `json:"prefix,omitempty"`
	Exact   *string `json:"exact,omitempty"`
	Rewrite *bool   `json:"rewrite,omitempty"`
}

func ConfigFromRequest(req *commonHttp.Request) (*Config, error) {
	conf := &Config{}
	return conf, json.Unmarshal([]byte(req.KulyData.Step.Config), &conf)
}
