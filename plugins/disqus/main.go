package disqus

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/plugins"
)

func init() {
	plugins.Register("disqus", New)
}

func New(opts map[string]string, store plugins.FileAttrStore, auth plugins.AuthManager) plugins.Plugin {
	logrus.Debugf("init disqus, %+v", opts)
	return &Disqus{
		shortName: opts["short-name"],
	}
}

type Disqus struct {
	shortName string
}

func (d *Disqus) Footer(path string) (template.HTML, error) {
	if d.shortName == "" {
		return "", fmt.Errorf("shortName is not set")
	}
	return template.HTML(strings.Replace(tmpl, "%disqus_shortname%", d.shortName, 1)), nil
}

var tmpl = `
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
`
