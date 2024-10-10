package plugins

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/pubsub"
)

type PluginState struct {
	Name      string `gorm:"primaryKey"`
	IsEnabled bool
	Config    string
}

type IManager interface {
	// set plugin
	Enable(ctx context.Context, name string) error
	Disable(ctx context.Context, name string) error
	SetConfig(ctx context.Context, name string, configString string) error // will it re-initialize? or just pass to plugin

	// hooks for cluster
	TriggerFileReadHooks(path string, data []byte) ([]byte, error)
	GetFooters(path string) ([]template.HTML, error)
	HandleHTTPRequest(c *gin.Context)
}

type PluginInit func(ctx context.Context, conf any, store metadata.IStore, hub *pubsub.Hub) (Plugin, error)

var _ IManager = (*Manager)(nil)

type Manager struct {
	rootCtx context.Context

	db    *gorm.DB
	store metadata.IStore
	hub   *pubsub.Hub

	// plugins
	Plugins map[string]Plugin

	Footers   map[string]PluginFooter
	ReadHooks map[string]PluginReadHook
	Handlers  map[string]http.Handler
}

func New(ctx context.Context, db *gorm.DB, store metadata.IStore, hub *pubsub.Hub) (*Manager, error) {
	if err := db.AutoMigrate(&PluginState{}); err != nil {
		return nil, errors.WithStack(err)
	}

	// First, ensure all plugin(even disabled) have PluginState.
	for name, driver := range drivers {
		logrus.Tracef("make plugin entry: %s", name)

		if err := db.WithContext(ctx).FirstOrCreate(&PluginState{
			Name:   name,
			Config: driver.Default,
		}).Error; err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// Initialize enabled Plugin
	states := []PluginState{}
	if err := db.WithContext(ctx).Where(PluginState{IsEnabled: true}).Find(&states).Error; err != nil {
		return nil, err
	}

	manager := &Manager{
		rootCtx: ctx,

		db:    db,
		store: store,
		hub:   hub,

		Plugins:   map[string]Plugin{},
		Footers:   map[string]PluginFooter{},
		ReadHooks: map[string]PluginReadHook{},
		Handlers:  map[string]http.Handler{},
	}
	for _, state := range states {
		if err := manager.initialzePlugin(state.Name, state.Config); err != nil {
			return nil, err
		}
	}
	return manager, nil
}

func (manager *Manager) initialzePlugin(name string, config string) error {
	driver, err := getDriver(name)
	if err != nil {
		return err
		// TODO removed plugin, show it but not enabled & missing.
	}

	// load config
	conf := driver.newConfig()

	if err := yaml.Unmarshal([]byte(config), conf); err != nil {
		return errors.WithStack(err)
	}

	plugin, err := driver.Init(manager.rootCtx, conf, manager.store, manager.hub)
	if err != nil {
		return err
	}

	manager.Plugins[name] = plugin

	if v, ok := plugin.(PluginReadHook); ok {
		manager.ReadHooks[name] = v
		logrus.Tracef("read hook detected")
	}
	if v, ok := plugin.(PluginFooter); ok {
		manager.Footers[name] = v
		logrus.Tracef("footer detected")
	}
	if v, ok := plugin.(PluginHTTPHandler); ok {
		manager.Handlers[name] = v
		logrus.Tracef("route detected")
	}

	return nil
}

func (m *Manager) TriggerFileReadHooks(path string, data []byte) ([]byte, error) {
	return data, nil
}
func (m *Manager) GetFooters(path string) ([]template.HTML, error) {
	result := []template.HTML{}
	for _, hook := range m.Footers {
		data, err := hook.Footer(path)
		if err != nil {
			logrus.Error(err)
			result = append(result, template.HTML("error in plugin: "+err.Error()))
			//return result, err
		}
		result = append(result, template.HTML(data))
	}
	return result, nil
}

func (m *Manager) HandleHTTPRequest(c *gin.Context) {
	name := c.Param("name")

	if h, exist := m.Handlers[name]; exist {
		// TODO make header
		h.ServeHTTP(c.Writer, c.Request)
	}
}
func (m *Manager) List(ctx context.Context) ([]PluginState, error) {
	plugins := []PluginState{}
	if err := m.db.WithContext(ctx).Find(&plugins).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return plugins, nil
}
func (m *Manager) Get(ctx context.Context, name string) (*PluginState, error) {
	plugin := &PluginState{}
	if err := m.db.WithContext(ctx).Where(&PluginState{Name: name}).Take(plugin).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return plugin, nil
}

func (m *Manager) Enable(ctx context.Context, name string) error {
	state := &PluginState{}
	if err := m.db.WithContext(ctx).Where(&PluginState{
		Name: name,
	}).Take(state).Error; err != nil {
		return errors.WithStack(err)
	}
	state.IsEnabled = true

	if err := m.db.WithContext(ctx).Save(state).Error; err != nil {
		return errors.WithStack(err)
	}

	// initialize plugin
	if err := m.initialzePlugin(name, state.Config); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Disable(ctx context.Context, name string) error {
	state := &PluginState{}
	if err := m.db.WithContext(ctx).Where(&PluginState{
		Name: name,
	}).Take(state).Error; err != nil {
		return errors.WithStack(err)
	}
	state.IsEnabled = false

	if err := m.db.WithContext(ctx).Save(state).Error; err != nil {
		return errors.WithStack(err)
	}

	// remove plugin instance from manager
	delete(m.Plugins, name)
	delete(m.Footers, name)
	delete(m.ReadHooks, name)
	delete(m.Handlers, name)

	return nil
}
func (m *Manager) SetConfig(ctx context.Context, name string, configString string) error {
	state := PluginState{
		Name: name,
	}
	if err := m.db.Take(&state).Error; err != nil {
		return errors.WithStack(err)
	}

	state.Config = configString

	if err := m.db.Save(&state).Error; err != nil {
		return errors.WithStack(err)
	}

	if !state.IsEnabled {
		return nil
		// there is no plugin instance
	}

	drivers, err := getDriver(name)
	if err != nil {
		return err
	}
	conf := drivers.newConfig()
	if err := yaml.Unmarshal([]byte(configString), conf); err != nil {
		return err
	}

	plugin, exist := m.Plugins[name]
	if !exist {
		return errors.Errorf("plugin instance not found: %s", name)
	}
	return plugin.SetConfig(ctx, conf)
}
