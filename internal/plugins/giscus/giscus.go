package utteranc

import (
	"bytes"
	"html/template"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
	"github.com/bluemir/wikinote/internal/pubsub"
)

// https://giscus.app

type Giscus struct {
	*Options
	html string
}
type Options struct {
	Repo          string `yaml:"repo"`
	RepoId        string `yaml:"repo-id"`
	Category      string `yaml:"category"`
	CategoryId    string `yaml:"category-id"`
	Reactions     bool
	InputPosition string `yaml:"input-position"`
	Theme         string
	Lang          string
}

var defaultConfig string = `

repo:
repoId:
category:
categoryId:
reactions:
inputPosition: 
theme:
lang:
`

func init() {
	plugins.Register("giscus", New, defaultConfig, &Options{})
}

func New(ctx context.Context, conf any, store metadata.IStore, hub *pubsub.Hub) (plugins.Plugin, error) {
	opt := &Options{}
	opt, ok := conf.(*Options)
	if !ok {
		return nil, errors.Errorf("option type not matched: %T", conf)
	}

	// validation
	if opt.Repo == "" {
		return nil, errors.Errorf("opt.Repo required")
	}
	if opt.RepoId == "" {
		return nil, errors.Errorf("opt.RepoId required")
	}
	if opt.Category == "" {
		return nil, errors.Errorf("opt.Category required")
	}
	if opt.CategoryId == "" {
		return nil, errors.Errorf("opt.Category required")
	}

	// default values
	if opt.InputPosition == "" {
		opt.InputPosition = "bottom"
	}
	if opt.Theme == "" {
		opt.Theme = "light"
	}
	if opt.Lang == "" {
		opt.Lang = "en"
	}

	html, err := makeHTML(opt)
	if err != nil {
		return nil, err
	}

	return &Giscus{opt, html}, nil
}
func (giscus *Giscus) Footer(path string) ([]byte, error) {
	return []byte(giscus.html), nil
}

var html = `
<style>
	.giscus, .giscus-frame {
		width: 100%;
		border: none;
	}
</style>
<script src="https://giscus.app/client.js"
	data-repo="{{ .Repo }}"
	data-repo-id="{{ .RepoId }}"
	data-category="{{ .Category }}"
	data-category-id="{{ .CategoryId }}"
	data-mapping="pathname"
	data-reactions-enabled="{{ if .Reactions }}1{{else}}0{{end}}"
	data-emit-metadata="0"
	data-inputPosition="{{ .InputPosition }}"
	data-theme="{{ .Theme }}"
	data-lang="{{ .Lang }}"
	data-loading="lazy"
	>
</script>
`

func makeHTML(opt *Options) (string, error) {
	buf := bytes.NewBuffer(nil)
	t, err := template.New("").Parse(html)
	if err != nil {
		return "", err
	}
	if err := t.Execute(buf, opt); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (*Giscus) SetConfig(ctx context.Context, conf any) error {
	_, ok := conf.(*Options)
	if !ok {
		return errors.Errorf("optiontype not matched: %T", conf)
	}
	return nil
}
