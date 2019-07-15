//minimal lib
var $ = {
	get: function(target, query) {
		if(typeof target.querySelector !== "function") {return $.get(document, target)}
		return target.querySelector(query);
	},
	all: function(target, query) {
		if(typeof target.querySelectorAll !== "function") {return $.all(document, target)}
		return Array.prototype.slice.call(target.querySelectorAll(query));
	},
	create: function(tagname, attr) {
		var newTag = document.createElement(tagname);
		if (attr && attr.$text){
			newTag.appendChild(document.createTextNode(attr.$text));
		}
		if (attr && attr.$html){
			newTag.innerHTML = attr.$html;
		}
		if (attr && attr.$child) {
			newTag.appendChild(attr.$child)
		}
		for(var key in (attr || {})){
			if (key[0] == "$") {
				continue; //skip
			}
			newTag.setAttribute(key, attr[key]);
		}
		return newTag;
	},
	request: async function $request(method, url, options) {
		var opts = options || {};

		if (opts.timestamp !== false) {
			opts.query = opts.query || {};
			opts.query["_timestamp"] = Date.now();
		}

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

			if (opts.auth) {
				console.debug("request with auth", opts.auth)
				// In Chrome and firefox Auth heaer not included request(due to security, see https://bugs.chromium.org/p/chromium/issues/detail?id=128323)
				// so forced set header
				req.open(method, resolveParam(url, opts.params) + queryString(opts.query), true, opts.auth.user, opts.auth.password);
				req.setRequestHeader("Authorization", "Basic " + btoa(opts.auth.user+":"+opts.auth.password));
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
	},
	timeout: async function(ms) {
		return new Promise(function(resolve, reject){
			setTimeout(resolve, ms);
		});
	},
	defer: function() {
		var ret = {}
		ret.promise = new Promise(function(resolve, reject){
			ret.resolve = resolve;
			ret.reject = reject;
		});
		return ret;
	},
	prevent: function(func){
		return function(evt){
			evt.preventDefault();
			return func();
		}
	},
	template: function(strings, ...args) {
		var html = ""

		for( var i = 0; i < strings.length; i++) {
			html += strings[i] + (args[i] ||"");
		}

		var template = $.create("template", { $html: html });
		return template;
	},
	render: function(templateNode, data) {
		var clone = document.importNode(templateNode.content, true);

		var f = function(match, name) {
			var arr = name.split(".")
			var result = data;
			for (var i = 0; i < arr.length; i++) {
				if (!result) {
					console.warn(`[$.template] find pattern '${name}', but not provided`);
					return "";
				}
				result = result[arr[i]]
			}
			return result || "";
		}
		var pattern = /{{\s*([a-zA-Z0-9._-]+)\s*}}/g
		var each = function(node) {
			switch (node.nodeType) {
				case Node.TEXT_NODE:
					node.textContent = node.textContent.replace(pattern, f)
					break;
				case Node.ELEMENT_NODE:
					for (var i = 0; i <  node.attributes.length; i++) {
						node.attributes[i].value = node.attributes[i].value.replace(pattern, f);
					}
					break;
				default:
			}
			node.childNodes.forEach(each)
		}
		clone.childNodes.forEach(each)

		return clone;
	},
	form: function(form) {
		return $.all(form, "input").map((e) => {
			var name = e.attr("name");
			var value = e.value;
			var t = e.attr("type");
			if (t == "number") {
				value = value - 0; // change to number
			}
			return {name, value}
		}).reduce((p, c) => {
			p[c.name] = c.value;
			return p;
		},{})
	},
	bindForm: function(form, data) {
		$.all(form, "input").forEach((e) => {
			var name = e.attr("name");
			e.value = data[name] || "";
		});
	},
	filters: {
		exceptTemplate: function(elem) {
			return elem.tagName != "TEMPLATE";
		}
	},
	event: new EventTarget(),
	_registerGlobal: function() {
		window.$ = this;
	},
}

function resolveParam(url, params) {
	if (params == null) {
		return url
	}
	return url.replace(/:([a-zA-Z0-9]+)/g, function(matched, name){
		if (params[name]) {
			return params[name];
		}
		console.warn(`[$.reqeust] find param pattern '${name}', but not provided`);
		return matched;
	});
}

function queryString(obj) {
	if (obj == null) {
		return "";
	}
	return "?" + Object.keys(obj).map(function(key) {
		return key + "=" + obj[key];
	}).join("&");
}

function extend(TargetClass, proto){
	if (TargetClass.hasOwnProperty("__minilib_inserted__")) {
		console.trace("already installed")
		return // already inserted
	}

	Object.keys(proto).forEach(function(name) {
		if (name  in TargetClass.prototype) {
			console.warn(`cannot extend prototype: '${name}' already exist`)
			return; // skip
		}
		TargetClass.prototype[name] = proto[name];
	});

	TargetClass.__minilib_inserted__ = true
}
extend(Node, {
	removeThis : function(){
		this.parentElement.removeChild(this);
	},
	clear : function(filter) {
		var f = filter || function(e) { return true };
		this.childNodes.filter(f).forEach((e) => this.removeChild(e))
	}
});

extend(Element, {
	attr: function(name, value){
		if (value === null) {
			this.removeAttribute(name);
			return
		}
		if (value !== undefined) {

			this.setAttribute(name, value)
			return value;
		} else {
			return this.getAttribute(name)
		}
	},
})

extend(EventTarget, {
	on : function() {
		var listeners = [];
		var eventNames = [];
		for (var i = 0; i < arguments.length; i++) {
			switch(typeof(arguments[i])) {
				case "function":
					listeners.push(arguments[i]);
					break;
				case "string":
					eventNames.push(arguments[i]);
					break;
				default:
					throw Error("'on' only accept function or string")
			}
		}
		eventNames.forEach((name) => {
			listeners.forEach((func) => {
				this.addEventListener(name, func)
			})
		})
		return this;
	},
	fireEvent: function(name, detail) {
		var evt = new CustomEvent(name, {detail: detail});
		this.dispatchEvent(evt);
		return this;
	}
});

extend(NodeList, {
	"map": Array.prototype.map,
	"filter": Array.prototype.filter,
	//"forEach": Array.prototype.forEach,
});

extend(Array, {
	"all": function all() { return Promise.all(this); },
	"race": function race() { return Promise.race(this); },

	// with the lovely addiction of ...
	"any": function any() { return Promise.any(this); },
});


class CustomElement extends HTMLElement {
	constructor(template) {
		super();
		if (arguments.length != 1 || !arguments[0]) {
			throw Error("CustomElement Must have 1 arguments");
		}
		if (!template instanceof DocumentFragment) {
			throw Error("not a DocumentFragment");
		}
		var clone = document.importNode(template, true)
		var shadow = this["--shadow"] = this.attachShadow({mode: 'open'})
		shadow.appendChild(clone);
	}

	attributeChangedCallback(name, oldValue, newValue) {
		//  to use set follow to custom elements
		//
		//	static get observedAttributes() {
		//		return ["cluster"];
		//	}
		this.fireEvent("attribute-changed", {
			name: name,
			old: oldValue,
			new: newValue,
		});
	}
	connectedCallback()  {
		this.fireEvent("connected")
	}
	disconnectedCallback() {
		this.fireEvent("disconnected")
	}
}

$.CustomElement = CustomElement;

export default $;
