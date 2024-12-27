import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<style>
		@import "/static/css/root.css";

		:host {
			display: inline-grid;
		}
		textarea {
			resize: none;
			overflow: hidden;
		}
	</style>
	<textarea @input="${evt => elem.resize(evt)}">${elem.value}</textarea>
`;

class Textarea extends $.CustomElement {
	static get formAssociated() {
		return true;
	}

	static get observedAttributes() {
		return ["name"];
	}

	#internal = this.attachInternals();
	#value

	get value(){
		return this.#value;
	}
	set value(v){
		this.#value = v;
		this.#internal.setFormValue(v);
	}

	constructor() {
		super();
	}
	onConnected() {
		this.resize();
	}
	formResetCallback() {
		$.get(this.shadowRoot, "textarea").value = "";
		this.value = "";
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}

	resize() {
		let $textarea = $.get(this.shadowRoot, "textarea");
		$textarea.style.height = `auto`;
		$textarea.style.height = `${$textarea.scrollHeight + 2}px`;
	}
}
customElements.define("c-textarea", Textarea);
