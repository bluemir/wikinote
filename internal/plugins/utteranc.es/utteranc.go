package utteranc

import (
	"bytes"
	"html/template"

	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/internal/plugins"
)

type Utteranc struct {
	*Options
	html string
}
type Options struct {
	Repo  string
	Label string
	Theme string
}

func init() {
	plugins.Register("utteranc", New, &Options{})
}

func New(o interface{}, store *plugins.Store) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	// validation
	if opt.Repo == "" {
		return nil, errors.Errorf("opt.Repo must required")
	}

	// default values
	if opt.Label == "" {
		opt.Label = "comments"
	}
	if opt.Theme == "" {
		opt.Theme = "github-light"
	}

	return &Utteranc{opt, makeHTML(opt)}, nil
}
func (utteranc *Utteranc) Footer(path string) ([]byte, error) {
	return []byte(utteranc.html), nil
}

var html = `
<script src="https://utteranc.es/client.js"
        repo="{{.Repo}}"
        issue-term="pathname"
        label="{{.Label}}"
        theme="{{.Theme}}"
        crossorigin="anonymous"
        async>
</script>
`

func makeHTML(opt *Options) string {
	buf := bytes.NewBuffer(nil)
	template.Must(template.New("").Parse(html)).Execute(buf, opt)
	return buf.String()
}
