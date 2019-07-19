package authz

import (
	"github.com/mgutz/str"

	"github.com/bluemir/wikinote/plugins"
)

type Rule struct {
	Object  Match    `yaml:"object"`
	Subject Match    `yaml:"subject"`
	Actions []string `yaml:"action"`
}
type Match struct {
	All []KV `yaml:"all"`
	Any []KV `yaml:"any"`
}
type KV struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
type Attr interface {
	Attr(key string) string
}

func (rule *Rule) match(ctx *plugins.AuthContext) bool {
	// Look action
	if !str.SliceContains(rule.Actions, ctx.Action) {
		return false
	}
	// Match any
	if matchAny(ctx.Object, rule.Object.Any) != true {
		return false
	}
	if matchAny(ctx.Subject, rule.Subject.Any) != true {
		return false
	}
	// Match all
	if matchAll(ctx.Object, rule.Object.All) != true {
		return false
	}
	if matchAll(ctx.Subject, rule.Subject.All) != true {
		return false
	}
	return true
}
func matchAny(attr Attr, kvs []KV) bool {
	if len(kvs) == 0 {
		return true
	}
	for _, kv := range kvs {
		if matchSingle(attr, kv) {
			return true
		}
	}
	return false
}
func matchAll(attr Attr, kvs []KV) bool {
	for _, kv := range kvs {
		if !matchSingle(attr, kv) {
			return false
		}
	}
	return true
}
func matchSingle(attr Attr, kv KV) bool {
	/*
		logrus.
			WithField("method", "plugins.authz.matchSingle").
			Tracef("attr[%s]=%s, v=%s", kv.Key, attr.Attr(kv.Key), kv.Value)
	*/
	if kv.Value == "" {
		//just check exist
		return attr.Attr(kv.Key) != kv.Value
	}
	return attr.Attr(kv.Key) == kv.Value
}
