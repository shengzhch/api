package util

import (
	"github.com/bitly/go-simplejson"
)

type Config struct {
	vals *simplejson.Json
}

func NewConfig() (*Config, error) {
	//todo
	return &Config{vals: simplejson.New()}, nil
}

func (this *Config) Set(key string, v interface{}) {
	this.vals.Set(key, v)
}

func (this *Config) Get(key string) *simplejson.Json {
	return this.vals.Get(key)
}

func (this *Config) GetMuststring(key string) string {
	return this.vals.Get(key).MustString()
}

func (this *Config) Map() (map[string]interface{}, error) {
	return this.vals.Map()
}

func (this *Config) Add(js *simplejson.Json) {
}

func (this *Config) AddConfig(c *Config) {
	cmap, _ := c.Map()
	for k, v := range cmap {
		this.Set(k, v)
	}
}
func (this *Config) Getvals() *simplejson.Json {
	return this.vals
}
