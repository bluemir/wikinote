package discus

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/pkg/plugins"
)

type Options struct {
	ShortName string `yaml:"shortName"`
}
type Discus struct {
	*Options
	html string
}

func init() {
	plugins.Register("discus", New, &Options{})
}

func New(o interface{}, store *plugins.Store) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	return &Discus{opt, makeHTML(opt.ShortName)}, nil
}
func (discus *Discus) Footer(path string) ([]byte, error) {
	return []byte(discus.html), nil
}
func makeHTML(shortName string) string {
	return strings.Replace(`
<div id="disqus_thread"></div>
<script type="text/javascript">
/* * * CONFIGURATION VARIABLES: EDIT BEFORE PASTING INTO YOUR WEBPAGE * * */
var disqus_shortname = '%disqus_shortname%'; // required: replace example with your forum shortname

/* * * DON'T EDIT BELOW THIS LINE * * */
(function() {
	var dsq = document.createElement('script'); dsq.type = 'text/javascript'; dsq.async = true;
	dsq.src = '//'+disqus_shortname+'.disqus.com/embed.js';
	(document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
})();
</script>
<noscript>Please enable JavaScript to view the <a href="https://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>
`, "%disqus_shortname%", shortName, 1)

}
