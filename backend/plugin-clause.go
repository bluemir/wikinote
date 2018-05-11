package backend

import "html/template"

type PluginClause interface {
	Footer(path string) []PluginResult
	AfterWikiSave(path string, data []byte) error
}

type pluginClause backend

func (b *pluginClause) Footer(path string) []PluginResult {
	result := []PluginResult{}
	for _, p := range b.plugins.footer {
		d, e := p.Footer(path)
		result = append(result, PluginResult{
			d, e,
		})
	}
	return result
}
func (b *pluginClause) AfterWikiSave(path string, data []byte) error {
	for _, p := range b.plugins.afterWikiSave {
		e := p.AfterWikiSave(path, data)
		if e != nil {
			return e
		}
	}
	return nil
}

type PluginResult struct {
	Data template.HTML
	Err  error
}
