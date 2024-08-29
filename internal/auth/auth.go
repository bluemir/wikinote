package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/bluemir/wikinote/internal/initializer"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

// config를 변경할때는 DB에 저장하고 Memory에도 반영한다.
// load 할때는 DB에서 가져 온다.
// 패턴화 할수 있을듯

type IManager interface {
	// Authn
	Default(username, unhashedKey string) (*User, error)
	HTTP(req *http.Request) (*User, error)

	// Authz
	Can(user *User, verb Verb, resource Resource) error

	// User
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, name string) (*User, bool, error)
	ListUsers(ctx context.Context) ([]User, error)
	//UpdateUser(ctx context.Context) (*User, error)
	DeleteUser(ctx context.Context, name string) error

	// Token
	IssueToken(ctx context.Context, username, unhashedKey string, opts ...TokenOpt) (*Token, error)
	GenerateToken(ctx context.Context, username string, opts ...TokenOpt) (*Token, string, error)
	//GetToken(ctx context.Context, username, unhashedKey string) (*Token, error)
	//ListToken(ctx context.Context, username string) ([]Token, error)
	//RevokeToken(ctx context.Context, username, unhashedKey string) error

	// Role
	CreateRole(ctx context.Context, name string, rule []Rule) (*Role, error)
	GetRole(ctx context.Context, name string) (*Role, error)
	ListRoles(ctx context.Context) ([]Role, error)
	//UpdateRole(ctx context.Context, role *Role) (*Role, error)
	DeleteRole(ctx context.Context, name string) error

	// Assign
	AssignRole(ctx context.Context, subject Subject, roles ...string) error
	DiscardRole(ctx context.Context, subject Subject, roles ...string) error
}

var _ IManager = &Manager{}

type Manager struct {
	conf *Config
	db   *gorm.DB
}
type Config struct {
	Salt  string
	Group struct {
		NewUserGroups []string
	}
}

func New(ctx context.Context, db *gorm.DB) (*Manager, error) {
	if err := db.AutoMigrate(
		&User{},
		&Group{},
		&Role{},
		&Token{},
		&Assign{},
	); err != nil {
		return nil, errors.WithStack(err)
	}
	// check initialized
	// if not, manager

	conf := Config{
		/* Initial Config */
		Salt: hash(
			xid.New().String(),
			"wikinote",
			time.Now().String(),
		),
		Group: struct {
			NewUserGroups []string
		}{
			NewUserGroups: []string{"user"},
		},
	}

	if err := initializer.LoadOrInit(
		ctx, db, "auth", &conf,
		//initializeConfig(db),
		initializeDefaultObject(db),
	); err != nil {
		return nil, err
	}

	/*
		store := NewConfigStore()
		if err := store.Load(ctx, "auth", AuthConfig{}, InitializeFunc); err != nil {
			return err
		}
		if err := store.Save(ctx, "auth", AuthConfig{}); err != nil {
			return err
		}
	*/

	return &Manager{
		conf: &conf,
		db:   db,
	}, nil
}
