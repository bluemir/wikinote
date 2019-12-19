import $ from "../lib/minilib.module.js";

var template = $.template`
<style>
:host {
}
input {
width: 300px;
}
</style>
<form>
	<section>
	</section>
	<button>+</button>
</form>
`
var kvTemplate = $.template`
<p>
	<input name="key" value="{{.key}}" />
	<input name="value" value="{{.value}}" />
	<button x-func="delete">-</button>
</p>
`
class KeyValueEditor extends $.CustomElement {
	constructor() {
		super(template.content);
		$.get(this.shadow, "button").on("click", e => this.add(e))
		this.$editor.on("click", e => {
			// e.path[i].matches("p")
			// or
			e.preventDefault();
			if (!e.path[0].matches("button[x-func=delete]")) {
				return
			}
			var p = e.path[0].closest("p")
			if(!p) {
				return
			}

			p.removeThis();
		})
	}
	add(e) {
		this.$editor.appendChild($.render(kvTemplate, {}))
		e.preventDefault();
	}
	get $editor() {
		return $.get(this.shadow, "section");
	}
	set data(v) {
		this.$editor.clear();
		Object.entries(v).map(([key, value]) => {
			 return $.render(kvTemplate, {key, value})
		}).reduce($.util.reduce.appendChild, this.$editor);
	}
	get data() {
		console.log($.all(this.$editor, "p"))
		return $.all(this.$editor, "p").map(elem => {
			return {
				key:   $.get(elem, "input[name=key]").value,
				value: $.get(elem, "input[name=value]").value,
			}
		}).reduce((obj, kv) => {
			obj[kv.key] = kv.value;
			return obj;
		}, {});
	}
}
customElements.define("kv-editor", KeyValueEditor);
