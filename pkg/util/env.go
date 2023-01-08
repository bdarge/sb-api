package util

import (
	"os"
	"strconv"
)

type EnvGetter struct{}

type ConfigSource interface {
	GetString(name string) string
}

type Env struct {
	source ConfigSource
}

func (r *EnvGetter) GetString(name string) string {
	return os.Getenv(name)
}
func NewEnv() *Env {
	return &Env{source: &EnvGetter{}}
}

func (c *Env) GetBool(name string) bool {
	s := c.source.GetString(name)
	i, err := strconv.ParseBool(s)
	if nil != err {
		return false
	}
	return i
}
