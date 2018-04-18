package user

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend/config"
)

type Manager interface {
	Get(username string) (*User, error)
	Put(user *User) error
	List() ([]User, error)
	New(user *User, token string) error
	Delete(username string) error
	Auth(username string, password string) bool
}

func NewManager(db *gorm.DB, conf *config.Config) (Manager, error) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Token{})
	return &manager{
		db:   db,
		conf: conf,
	}, nil
}

type manager struct {
	db   *gorm.DB
	conf *config.Config
}

func (m *manager) key(id string) string {
	return fmt.Sprintf("/user/%s", id)
}

func (m *manager) Get(id string) (*User, error) {
	logrus.Debugf("find %s db", id)
	user := &User{}
	m.db.Where("name = ?", id).Take(user)

	return user, nil
}
func (m *manager) List() ([]User, error) {

	users := []User{}
	result := m.db.Find(users)

	return users, result.Error
}
func (m *manager) New(user *User, token string) error {
	result := m.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	logrus.Info(user)

	result = m.db.Create(&Token{
		UserID:    user.ID,
		HashedKey: hash(user.Name, token),
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (m *manager) Put(user *User) error {
	return m.db.Save(user).Error
}
func (m *manager) Delete(username string) error {
	return m.db.Where("name = ?", username).Delete(&User{}).Error
}
func (m *manager) Auth(username string, password string) bool {
	cnt := 0
	result := m.db.Table("tokens").
		Joins("JOIN users ON users.id = tokens.user_id").
		Where("users.name = ? AND tokens.hashed_key = ?", username, hash(username, password)).Count(&cnt)
	if result.Error != nil {
		// TODO maybe warning...
		return false
	}
	return cnt > 0
}
