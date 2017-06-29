package user

import (
	"encoding/json"
	"fmt"

	"github.com/docker/libkv/store"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend/config"
)

type Manager interface {
	Get(username string) (*User, error)
	Put(user *User) error
	List() ([]User, error)
	New(user *User) error
	Delete(username string) error
}

func NewManager(kv store.Store, conf *config.Config) (Manager, error) {
	return &manager{
		kv:   kv,
		conf: conf,
	}, nil
}

type manager struct {
	kv   store.Store
	conf *config.Config
}

func (m *manager) key(id string) string {
	return fmt.Sprintf("/user/%s", id)
}

func (m *manager) Get(id string) (*User, error) {
	logrus.Debugf("find %s on kv store", id)
	user := &User{}
	pair, err := m.kv.Get(m.key(id))
	if err != nil {
		return nil, err
	}

	logrus.Debugf("[USER GET] %s - %s:", pair.Key, pair.Value)
	err = json.Unmarshal(pair.Value, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (m *manager) List() ([]User, error) {
	v, err := m.kv.List("/user")
	if err != nil {
		if isNotFound(err) {
			return []User{}, nil
		}
		return nil, err
	}
	result := []User{}
	for _, pair := range v {
		u := &User{}

		err := json.Unmarshal(pair.Value, u)
		if err != nil {
			return nil, err
		}

		result = append(result, *u)
	}
	return result, nil
}
func (m *manager) New(user *User) error {
	user.Role = m.conf.User.Default.Role

	u, err := m.Get(user.Id)
	if err != nil && !isNotFound(err) {
		return err
	}
	if u != nil {
		return fmt.Errorf("user already exist: %s", u.Id)
	}

	buf, err := json.Marshal(user)
	if err != nil {
		return err
	}
	logrus.Debugf("save user %s", buf)
	return m.kv.Put(m.key(user.Id), buf, nil)
}
func (m *manager) Put(user *User) error {
	if buf, err := json.Marshal(user); err != nil {
		return err
	} else {
		return m.kv.Put(m.key(user.Id), buf, nil)
	}
}
func (m *manager) Delete(username string) error {
	return m.kv.Delete(m.key(username))
}
func isNotFound(err error) bool {
	return err == store.ErrKeyNotFound
}
