package auth

import (
	"context"
	_ "embed"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

//go:embed init_default_policy.yaml
var defaultPolicy []byte

func initializeDefaultObject(db *gorm.DB) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		txn := db.WithContext(ctx).Begin()
		defer txn.Rollback()

		data := struct {
			Roles   []Role
			Assigns []struct {
				Subject Subject
				Roles   []string
			}
		}{}

		if err := yaml.Unmarshal(defaultPolicy, &data); err != nil {
			return errors.WithStack(err)
		}

		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			buf, _ := yaml.Marshal(data)
			logrus.Debugf("%s", string(buf))
		}

		for _, role := range data.Roles {
			if err := txn.Create(&role).Error; err != nil {
				return errors.WithStack(err)
			}
		}
		for _, assign := range data.Assigns {
			if err := txn.Create(&Assign{
				Subject: assign.Subject,
				Roles:   setFromArray(assign.Roles),
			}).Error; err != nil {
				return errors.WithStack(err)
			}
		}

		return txn.Commit().Error
	}
}

func loadConfigFromDB(ctx context.Context, db *gorm.DB) (*Config, error) {
	conf := Config{}
	if err := db.WithContext(ctx).Take(&conf).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return &conf, nil
}
