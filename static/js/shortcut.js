import * as $ from "../lib/minilib.module.js";


// tab + space
// alt + s
// ctrl + s
//
const noop = ()=>{}

//var L = console;
var L = {info: noop, log: noop}

class Shortcut {
	constructor(elem) {
		if (!(elem instanceof EventTarget)) {
			console.debug(elem)
			throw Error("type error")
		}

		this.handlers = [];
		this.status = {};

		elem.on("keydown", evt => this.handleKeydown(evt));
		elem.on("keyup"  , evt => this.handleKeyup  (evt));
	}
	add(q, func) {
		var code = new Code(q, func, {preventDefault:true});
		this.handlers.push(code);
	}
	hook(q, func) {
		var code = new Code(q, func, {preventDefault:false});
		this.handlers.push(code);
	}
	handleKeydown(evt) {
		if (evt.keyCode == 229) {
			return // just ignore this event
			// http://lists.w3.org/Archives/Public/www-dom/2010JulSep/att-0182/keyCode-spec.html
			// it is 'in progress key input' code
		}
		this.status[evt.keyCode] = true;

		this.handleModifier(evt);

		this.fire(evt);
	}
	handleKeyup(evt) {
		delete this.status[evt.keyCode];

		this.handleModifier(evt);
	}
	handleModifier(evt) {
		if (evt.altKey) {
			this.status[toCode.alt] = true;
		} else {
			delete(this.status[toCode.alt])
		}

		if (evt.metaKey) {
			this.status[toCode.meta] = true;
		} else {
			delete(this.status[toCode.meta])
		}

		if (evt.shiftKey) {
			this.status[toCode.shift] = true;
		} else {
			delete(this.status[toCode.shift])
		}

		if (evt.ctrlKey) {
			this.status[toCode.ctrl] = true;
		} else {
			delete(this.status[toCode.ctrl])
		}
	}
	fire(evt) {
		var keyCodes = Object.keys(this.status).map(c => c-0).sort();
		L.log(keyCodes);

		this.handlers.forEach((h) => {
			if (h.match(keyCodes)) {
				h.func();
				if (h.preventDefault) {
					evt.preventDefault();
				}
			}
		});
	}
}

class Code {
	constructor(q, func, {preventDefault = true} = {}) {
		this.parse(q)
		this.func = func;
		this.preventDefault = preventDefault;
	}
	parse(q) {
		var keys = q.toLowerCase().split("+").map(s => s.trim());
		var codes = keys.map(c => toCode[c]).reduce((arr, code) => {
			arr.push(code);
			return arr
		}, []).sort();

		L.log(codes)
		this.codes = codes;
	}
	match(codes) {
		L.log(JSON.stringify(codes) , JSON.stringify(this.codes))
		return JSON.stringify(codes) == JSON.stringify(this.codes); // simple match
	}
}

const toCode = {
	"alt": 18, "ctrl": 17, "shift": 16,

	'esc':27, 'escape':27,
	'tab':9,
	'space':32,
	'return':13, 'enter':13,
	'backspace':8,

	'scrolllock':145, 'scroll_lock':145, 'scroll':145,
	'capslock':20, 'caps_lock':20, 'caps':20,
	'numlock':144, 'num_lock':144, 'num':144,

	'pause':19, 'break':19,
	'insert':45, 'home':36, 'delete':46, 'end':35,

	'pageup':33, 'page_up':33,
	'pagedown':34, 'page_down':34,

	'left':37, 'up':38, 'right':39, 'down':40,

	'f1':112, 'f2':113, 'f3':114, 'f4':115, 'f5':116,
	'f6':117, 'f7':118, 'f8':119, 'f9':120, 'f10':121,
	'f11':122, 'f12':123
};

(function init(){
	var str = "`1234567890-=qwertyuiop[]\\asdfghjkl;'zxcvbnm,./";
	for (var i= 0; i < str.length; i++) {
		var char = str.charAt(i);
		var code = char.toUpperCase().charCodeAt(0);

		toCode[char] = code;
	}
})();

export {Shortcut};
