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
		${elem.items.map((item, index) => html`
			<div>
				<input
					.value="${item}"
					@change="${evt => elem.setValue(index, evt.target.value)}"
					placeholder="${elem.attr("placeholder")}"
					type="${elem.attr("type")||"text"}"
				/>
				<c-button><button  @click="${evt => elem.removeValue(index)}">Delete</button></c-button>
			</div>
		`)}
		<c-button><button @click="${evt => elem.addValue(evt)}">Add</button></c-button>
	</div>
`;

// key value input
// encode to `{item},{item}`
class CustomElement extends $.CustomElement {
	static get formAssociated() {
		return true;
	}

	constructor() {
		super();
	}

	items = [""];
	#internal = this.attachInternals();

	setValue(index, value) {
		this.items[index] = value;
		this.#internal.setFormValue(this.items.join(","));
	}
	addValue() {
		this.items.push("");
		this.render();
	}
	removeValue(index) {
		this.items.splice(index, 1);
		this.render();
	}
	get value() {
		return this.items;
	}
	set value(v) {
		this.items = v;
		this.render();
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}
}
customElements.define("c-input-list", CustomElement);
