package backend

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type UserClause interface {
	Get(username string) (*User, error)
	Put(user *User) error
	List() ([]User, error)
	New(user *User, token string) error
	Delete(username string) error
	Auth(username string, password string) (*User, bool, error)
}

type userClause backend

func (b *userClause) Get(username string) (*User, error) {
	logrus.Debugf("find %s db", username)
	user := &User{}
	b.db.Where("name = ?", username).Take(user)

	return user, nil
}
func (b *userClause) Put(user *User) error {
	return b.db.Save(user).Error
}
func (b *userClause) List() ([]User, error) {

	users := []User{}
	result := b.db.Find(users)

	return users, result.Error
}
func (b *userClause) New(user *User, token string) error {
	if user.Role == "" {
		user.Role = b.conf.User.Default.Role
	}

	result := b.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	logrus.Info(user)

	result = b.db.Create(&Token{
		UserID:    user.ID,
		HashedKey: hash(user.Name, token),
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (b *userClause) Delete(username string) error {
	return b.db.Where("name = ?", username).Delete(&User{}).Error
}
func (b *userClause) Auth(username string, password string) (*User, bool, error) {
	user := &User{}
	result := b.db.
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

type User struct {
	gorm.Model
	Name  string
	Email string
	Role  string
}
type Token struct {
	gorm.Model
	UserID    uint
	HashedKey string
}

func hash(username string, key string) string {
	salt := rawHashHex(key + "__salt__" + username)
	return rawHashBase64(salt[:64] + key + salt[64:])
}
func rawHashBase64(input string) string {
	h := crypto.SHA512.New()
	io.WriteString(h, input)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
func rawHashHex(str string) string {
	h := crypto.SHA512.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}
