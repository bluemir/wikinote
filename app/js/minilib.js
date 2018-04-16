//mininal lib
class $ {
	static get(target, query) {
		if(typeof target.querySelector !== "function") {return $.get(document, target)}
		return target.querySelector(query);
	}
	static all(target, query) {
		if(typeof target.querySelectorAll !== "function") {return $.all(document, target)}
		return Array.prototype.slice.call(target.querySelectorAll(query));
	}
	static create(tagname, attr) {
		var newTag = document.createElement(tagname);
		if (attr && attr.$text){
			newTag.appendChild(document.createTextNode(attr.$text));
		}
		if (attr && attr.$html){
			newTag.innerHTML = attr.$html;
		}
		for(var key in (attr || {})){
			if (key[0] == "$") {
				continue; //skip
			}
			newTag.setAttribute(key, attr[key]);
		}
		return newTag;
	}
	static async request(method, url, options) {
		var opts = options || {}

		return new Promise(function(resolve, reject) {
			var req = new XMLHttpRequest();

			Object.keys(opts.header || {}).forEach(function(name){
				req.setRequestHeader(name, opts.header[name]);
			});

			req.addEventListener("readystatechange", function(){
				if (req.readyState  == 4) {
					var result = {
						statusCode: req.status,
						text : req.responseText,
					};
					if (req.status >= 200, req.status < 300){
						if(req.getResponseHeader("Content-Type").includes("application/json")) {
							result.json = JSON.parse(result.text);
						}
						resolve(result)
					} else {
						if(req.getResponseHeader("Content-Type").includes("application/json")) {
							result.json = JSON.parse(result.text);
						}
						reject(result);
					}
				}
			});

			if (opts.$auth) {
				req.open(method, resolveParam(url, opts.params) + queryString(opts.query), true, opts.$auth.user, opts.$auth.password);
			} else {
				req.open(method, resolveParam(url, opts.params) + queryString(opts.query), true);
			}

			switch (typeof opts.body) {
				case "object":
					req.setRequestHeader("Content-Type", "application/json")
					req.send(JSON.stringify(opts.body))
					break;
				case "string":
					req.send(opts.body);
					break;
				case "undefined":
					req.send();
					break; // just skip
				default:
					reject("unknown type: req.body");
					break;
			}
		});
	}
	static async timeout(ms) {
		return new Promise(function(resolve, reject){
			setTimeout(resolve, ms);
		});
	}
	defer() {
		var ret = {}
		ret.promise = new Promise(function(resolve, reject){
			ret.resolve = resolve;
			ret.reject = reject;
		});
		return ret;
	}
	prevent(func){
		return function(evt){
			evt.preventDefault();
			return func();
		}
	}
}

function resolveParam(url, params) {
	if (params == null) {
		return url
	}
	return url.replace(/:([a-zA-Z0-9]+)/g, function(matched, name){
		if (params[name]) {
			return params[name];
		}
		console.warn("[$.reqeust] find param pattern '"+name+"', but not provided");
		return matched;
	});
}

function queryString(obj) {
	if (obj == null) {
		return "";
	}
	return "?" + obj.keys().map(function(key) {
		return key + "=" + obj[key];
	}).join("&");
}



function extend(TargetClass, proto){
	Object.keys(proto).forEach(function(name) {
		if (name  in TargetClass.prototype) {
			console.warn("cannot extend prototype: '"+name+"' already exist")
			return; // skip
		}
		TargetClass.prototype[name] = proto[name];
	});
}

extend(Element, {
	attr: function(name, value){
		if (value !== undefined) {
			this.setAttribute(name, value)
			return value;
		} else {
			return this.getAttribute(name)
		}
	},
	"removeThis" : function(){
		this.parentElement.removeChild(this);
	},
	"clear" : function(){
		while (this.childNodes.length > 0) {
			this.removeChild(this.childNodes[0]);
		}
	}
})

extend(EventTarget, {
	"on" : function() {
		this.addEventListener.apply(this, arguments);
		return this;
	}
});

extend(NodeList, {
	"map": Array.prototype.map,
	//"forEach": Array.prototype.forEach,
});
extend(Array, {
	"all": function all() { return Promise.all(this); },
	"race": function race() { return Promise.race(this); },

	// with the lovely addiction of ...
	"any": function any() { return Promise.any(this); },
});


// XXX preloading for custom elements
export default $;
