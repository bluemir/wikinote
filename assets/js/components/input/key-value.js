import * as $ from "bm.js/bm.module.js";
//import {html, render} from '/lib/lit-html/lit-html.js';
import {html, render} from 'lit-html';
import {css} from "common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
		}
	</style>
	<div>
		${elem.kvs.map((kv, index) => html`
			<div>
				<input .value="${kv.key}"   @change="${evt => elem.setKey(index, evt.target.value)}"   placeholder="key"   />
				<input .value="${kv.value}" @change="${evt => elem.setValue(index, evt.target.value)}" placeholder="value" />
				<c-button><button @click="${evt => elem.removeKeyValue(index)}">Delete</button></c-button>
			</div>
		`)}
		<c-button><button @click="${evt => elem.addKeyValue(evt)}">Add</button></c-button>
	</div>


`;

// key value input
// encode to `{key}={value},{key}={value}`
class CustomElement extends $.CustomElement {
	static get formAssociated() {
		return true;
	}

	static get observedAttributes() {
		return ["name", "label", "placeholder"];
	}

	kvs = [{key:"", value:""}];
	#internal = this.attachInternals();

	constructor() {
		super();
	}

	setKey(index, key) {
		this.kvs[index].key = key;
		this.encodeToInternal();
	}
	setValue(index, value) {
		this.kvs[index].value = value;
		this.encodeToInternal();
	}

	encodeToInternal() {
		let v = this.kvs.map(({key, value}) => `${key}=${value}`).join(",");
		this.#internal.setFormValue(v);
	}
	addKeyValue() {
		this.kvs.push({key:"", value:""});
		this.render();
	}
	removeKeyValue(index) {
		this.kvs.splice(index,1);
		this.render();
	}

	get value() {
		return this.kvs.reduce((obj, kv) => {
			obj[kv.key] = kv.value;
			return obj
		}, {})
	}
	set value(obj) {
		this.kvs = Object.entries(obj).map(([k, v]) => ({key: k, value: v}))
		this.render();
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}
}
customElements.define("c-input-key-value", CustomElement);
