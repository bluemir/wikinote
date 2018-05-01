package user

import (
	"fmt"
	"io/ioutil"
	"path"

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
	Auth(username string, password string) (*User, bool, error)
}

func NewManager(db *gorm.DB, conf *config.Config, wikipath string) (Manager, error) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Token{})

	root := &User{
		Name:  "root",
		Email: "root@wikinote",
		Role:  "root",
	}
	db.Where("name=?", "root").FirstOrCreate(root)
	key := RandomString(16)
	// always make new token. If forget root key? just restart it
	db.Where(&Token{UserID: root.ID}).Assign(&Token{HashedKey: hash("root", key)}).FirstOrCreate(&Token{})

	// Save to File
	// QUESTION or just print stdout?
	ioutil.WriteFile(path.Join(wikipath, ".app", ".root_token"), []byte(key), 0644)

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
func (m *manager) Auth(username string, password string) (*User, bool, error) {
	user := &User{}
	result := m.db.
		Joins("JOIN tokens ON users.id = tokens.user_id").
		Where("users.name = ? AND tokens.hashed_key = ?", username, hash(username, password)).First(user)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return user, true, nil
}
