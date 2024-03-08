import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
			display: inline-block;
			position: relative;
			padding-top: 0.5rem;
		}
		input {
			display: block;
			outline: none;
			height: 2rem;
			border: 0;
			border-bottom: 1px solid var(--gray-300);
		}
		input:focus, input:not(:placeholder-shown) {
			border-bottom: 1px solid var(--gray-800);
		}

		input::placeholder {
			opacity: 0;
		}
		input:focus::placeholder {
			opacity: 1;
		}
		label {
			display: block;
			position: absolute;
			top: 0.5rem;
		}
		input:focus + label, input:not(:placeholder-shown) + label {
			font-size: 0.5rem;
			top: 0;
		}
	</style>
	<input
		id="input" type="${elem.attr("type")}" name="${elem.name}" placeholder="${elem.attr("placeholder") || " " }"
		@input="${ evt => elem.value = evt.target.value}"
		@keypress="${ evt => elem.onKeyPress(evt)}"
	/>
	<label for="input">${elem.attr("label")}</label>
`;

// TODO
// - animation
// - proxy attribute of input
// - type of input(outline, underline, ...)

class Input extends $.CustomElement {
	static get formAssociated() {
		return true;
	}

	static get observedAttributes() {
		return ["name", "label", "placeholder", "type"];
	}

	#internal = this.attachInternals();
	#value

	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}

	get name(){
		return this.attr("name");
	}
	get value(){
		return this.#value;
	}
	set value(v){
		this.#value = v;
		this.#internal.setFormValue(v);
	}
	formResetCallback() {
		$.get(this.shadowRoot, "input").value = "";
		this.value = "";
	}
	onKeyPress(evt) {
		if (evt.code == "Enter" && this.#internal.form) {
			this.#internal.form.requestSubmit();
		}
	}
}
customElements.define("c-input", Input);
