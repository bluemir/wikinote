package auth

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func (manager *Manager) Authz(ctx *Context) Result {
	log := logrus.WithField("method", "auth.manager:authz")

	log.Tracef("user: %#v", ctx.Subject.User)
	log.Tracef("object: %#v", ctx.Object)
	log.Tracef("action: %s", ctx.Action)

	labels := map[string]string{"role/default": "true"}
	if ctx.Subject.User != nil {
		labels = ctx.Subject.User.Labels
	}
	log.Tracef("lable: %#v", labels)

	for k, _ := range labels {
		if !strings.HasPrefix(k, "role/") {
			continue
		}
		role := strings.TrimPrefix(k, "role/")

		// find role.. load from file when load auth module
		r := manager.roles[role]
		log.Tracef("role: %#v", r)

		for _, rule := range r.Rules {
			log.Trace("rule", rule)
			if rule.IsNotMatchedObject(ctx.Object) {
				continue
			}
			if rule.IsNotMatchedAction(ctx.Action) {
				continue
			}

			// all matched
			return Accept
		}
	}
	if ctx.Subject.User == nil {
		return NeedAuthn
	}

	return Reject
}
