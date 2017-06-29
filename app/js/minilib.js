//mininal lib
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
		for(var key in (attr || {})){
			if (key[0] == "$") {
				continue; //skip
			}
			newTag.setAttribute(key, attr[key]);
		}
		return newTag;
	},
	request: async function $request(method, url, options) {
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
						reject(result);
					}
				}
			});

			req.open(method, resolveParam(url, opts.params) + queryString(opts.query), true);

			switch (typeof opts.body) {
				case "object":
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

}

function resolveParam(url, params) {
	if (params == null) {
		return url
	}
	return url.replace(/:([a-zA-Z0-9]+)/g, function(matched, name){
		if (params[name]) {
			return params[name];
		}
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


var elementProto = {
	"remove" : function(){
		this.parentElement.removeChild(this);
	},
	"clear" : function(){
		while (this.childNodes.length > 0) {
			this.removeChild(this.childNodes[0]);
		}
	}
};
Object.keys(elementProto).forEach(function(name) {
	if (name  in Element.prototype) {
		return; // skip
	}
	Element.prototype[name] = elementProto[name];
});

var htmlElementProto = {
	"on" : function() {
		this.addEventListener.apply(this, arguments);
		return this;
	}
}
Object.keys(htmlElementProto).forEach(function(name) {
	if (name in HTMLElement.prototype) {
		return; // skip
	}
	HTMLElement.prototype[name] = htmlElementProto[name];
});
export default $;
