package auth

import (
	"context"
	"net/http"

	"github.com/bluemir/wikinote/internal/initializer"
	"github.com/pkg/errors"
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
	AssignRole(ctx context.Context, subject Subject, role Role) error
	DiscardRole(ctx context.Context, subject Subject, role Role) error
}

var _ IManager = &Manager{}

type Manager struct {
	conf *Config
	db   *gorm.DB
}
type Config struct {
	Salt  string
	Group struct {
		Unauthorized string
		Newcomer     []string
	} `gorm:"type:bytes;serializer:gob"`
}

func New(ctx context.Context, db *gorm.DB) (*Manager, error) {
	if err := db.AutoMigrate(
		&User{},
		&Group{},
		&Role{},
		&Token{},
		&Assign{},
		&Config{},
	); err != nil {
		return nil, errors.WithStack(err)
	}
	// check initialized
	// if not, manager

	if err := initializer.EnsureInitialize(
		ctx, db, "auth",
		initializeConfig(db),
		initializeDefaultRole(db),
	); err != nil {
		return nil, err
	}

	conf, err := loadConfigFromDB(ctx, db)
	if err != nil {
		return nil, err
	}

	return &Manager{
		conf: conf,
		db:   db,
	}, nil
}
