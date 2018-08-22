package backend

import (
	"fmt"
	"strings"

	"github.com/bluemir/go-utils/auth"
	"github.com/bluemir/wikinote/backend/config"
	"github.com/bluemir/wikinote/plugins"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	defaultRule = `
admin: view, edit, user
editor:  view, edit, attach
viewer:  view
guest:  view
	`
)

func UserInit(authMng auth.Manager) error {
	root, ok, err := authMng.GetUser("root")
	if err != nil {
		return err
	}
	if !ok {
		// it means first time

		logrus.Debug("put default rules")
		lines := strings.Split(defaultRule, "\n")
		for _, line := range lines {
			p := strings.SplitN(line, ":", 2)
			if len(p) != 2 {
				// skip error line
				continue
			}
			logrus.Debug(line)
			role := p[0]
			actions := strings.Split(p[1], ",")
			for _, action := range actions {
				logrus.Debugf("%s - %s", role, strings.Trim(action, " "))
				err := authMng.PutRule(role, strings.Trim(action, " "))
				if err != nil {
					return err
				}
			}
		}

		// create root user
		err := authMng.CreateUser(&auth.User{
			Name: "root",
			Role: "root",
			Labels: map[string]string{
				"email": "root@wikinote",
			},
		})
		if err != nil {
			return err
		}

		root, _, err = authMng.GetUser("root")
		if err != nil {
			return err
		}
	}

	// always root role.
	root.Role = "root"
	err = authMng.UpdateUser(root)
	if err != nil {
		return errors.Wrapf(err, "root user update failed")
	}

	// always make new token. If forget root key? just restart it
	key := RandomString(16)
	_, err = authMng.IssueToken("root", key)
	if err != nil {
		return err
	}

	// QUESTION save file or just print stdout?
	logrus.Infof("Root Token: %s", key)
	return nil
}

func dbInit(db *gorm.DB) error {
	/*
		// Auth
		if !db.HasTable(&Rule{}) {
			// only first time
			db.CreateTable(&Rule{})
			lines := strings.Split(defaultRule, "\n")
			for _, line := range lines {
				p := strings.SplitN(line, ":", 2)
				if len(p) < 2 {
					// skip error line
					continue
				}
				role := p[0]
				actions := strings.Split(p[1], ",")
				for _, action := range actions {
					rule := &Rule{Role: role, Action: strings.Trim(action, " ")}
					db.Where(rule).FirstOrCreate(rule)
				}
			}
		}

		// User & Token
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
		// QUESTION save file or just print stdout?
		logrus.Infof("Root Token: %s", key)
	*/
	return nil
}

type pluginList struct {
	footer         []plugins.FooterPlugin
	afterWikiSave  []plugins.AfterWikiSavePlugin
	registerRouter map[string]plugins.RegisterRouterPlugin
}

func loadPlugins(db *gorm.DB, conf *config.Config) (*pluginList, error) {
	// TODO can on/off
	pluginNames := plugins.List()
	pl := &pluginList{
		footer:         []plugins.FooterPlugin{},
		afterWikiSave:  []plugins.AfterWikiSavePlugin{},
		registerRouter: map[string]plugins.RegisterRouterPlugin{},
	}

	for _, name := range pluginNames {
		logrus.Debugf("pluginconf %+v, %s", conf.Plugins[name], name)
		pc, ok := conf.Plugins[name].(map[interface{}]interface{})
		if !ok {
			pc = map[interface{}]interface{}{}
		}

		p, err := plugins.New(name, db, flat(pc)) // TODO config
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		logrus.Debugf("plugin '%s' is initialize", name)
		if f, ok := p.(plugins.FooterPlugin); ok {
			pl.footer = append(pl.footer, f)
		}
		if a, ok := p.(plugins.AfterWikiSavePlugin); ok {
			pl.afterWikiSave = append(pl.afterWikiSave, a)
		}
		if a, ok := p.(plugins.RegisterRouterPlugin); ok {
			pl.registerRouter[name] = a
		}
	}
	return pl, nil
}
func flat(conf map[interface{}]interface{}) map[string]string {
	result := map[string]string{}
	for k, v := range conf {
		switch val := v.(type) {
		case map[interface{}]interface{}:
			res := flat(val)
			for resk, resv := range res {
				result[k.(string)+"."+resk] = resv
			}
		default:
			result[k.(string)] = fmt.Sprint(v)
		}
	}
	return result
}