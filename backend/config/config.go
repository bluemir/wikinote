package config

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	FrontPage  string `yaml:"front-page"`
	AutoBackup bool   `yaml:"auto-backup"`
	User       struct {
		Default struct {
			Role string `yaml:"role"`
		} `yaml:"default"`
	} `yaml:"user"`
	Plugins map[interface{}]interface{} `yaml:"plugins"`
}

func ParseConfig(path string) (*Config, error) {
	resources := rice.MustFindBox("../../resources")
	// Default Value

	conf := &Config{}
	yaml.Unmarshal(resources.MustBytes("default-config.yaml"), conf)

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Warn("config.yaml not exist.")
		return conf, nil
	}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
func (conf *Config) Save(path string) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
func (conf *Config) Set(path, value string) error {
	// TODO refactor this. for map
	keys := strings.Split(path, ".")

	switch keys[0] {
	case "plugins":
		// maybe use v.SetMapIndex
		cur := conf.Plugins
		for i, v := range keys[1:] {
			if i == len(keys)-2 {
				// last one
				cur[v] = value
				return nil
			}
			cur = cur[v].(map[interface{}]interface{})
		}
	default:
		v := reflect.ValueOf(conf)
		for _, key := range keys {
			if v.Kind() == reflect.Ptr {
				v = reflect.Indirect(v)
			}
			logrus.Debugf("key:%s, current Type: %s, %s", key, v.Type(), v)
			if v.Kind() != reflect.Struct {
				return fmt.Errorf("unexpected type, %s", v.Kind())
			}
			if i, ok := checkTag(v, "yaml", key); ok {
				v = v.Field(i)
			}
		}
		if !v.CanSet() {
			return fmt.Errorf("cannot set %s", path)
		}
		v.SetString(value)
	}
	return nil
}
func checkTag(v reflect.Value, tagName string, key string) (int, bool) {
	for i := 0; i < v.Type().NumField(); i++ {
		tag, ok := v.Type().Field(i).Tag.Lookup(tagName)
		if !ok {
			continue // next field
		}
		if tag == key {
			return i, true
		}
	}
	return -1, false
}
func (conf *Config) Get(path string) (interface{}, error) {
	keySlice := strings.Split(path, ".")
	v := reflect.ValueOf(conf)
	//iterate through field names ,ignore the first name as it might be the current instance name
	for _, key := range keySlice {
		if v.Kind() == reflect.Ptr {
			v = reflect.Indirect(v)
		}
		logrus.Debugf("key:%s, current Type: %s, %s", key, v.Type(), v)
		switch v.Kind() {
		case reflect.Struct:
			for i := 0; i < v.Type().NumField(); i++ {
				tag, ok := v.Type().Field(i).Tag.Lookup("yaml")
				if !ok {
					continue // next field
				}
				if tag == key {
					v = v.Field(i)
					break // next key
				}
			}
		case reflect.Map:
			v = v.MapIndex(reflect.ValueOf(key)).Elem()
		case reflect.Slice:
			if i, err := strconv.ParseInt(key, 10, 0); err != nil {
				return "", fmt.Errorf("%s is not int", key)
			} else {
				v = v.Index(int(i))
			}
		default:
			return "", fmt.Errorf("%s, expected []interface{} or map[string]interface{}; got %T", key, v)
		}
	}
	return v.Interface(), nil
}
