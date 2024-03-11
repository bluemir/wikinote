import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {

		}

		/* remove input default system */
		input {
			appearance: none;
			border: none;
			outline: none;
		}

		/* */

		div {
			position: relative;
			label {
				color: gray;
				position: absolute;
				font-size: 0.7rem;
				top: 0;
				left: 0;
			}
			input {
				border-bottom: 1px solid gray;
			}
			input::placeholder {
				opacity: 0;
			}

			padding-top: 0.7rem;
		}
		
		div:has(input:placeholder-shown) {
			label {
				font-size: 1rem;
				top: 0.7rem;
			}
		}
		div:has(input:not(:placeholder-shown)) {
			label {
				color: black;
			}
		}
		div:has(input:focus) {
			label {
				top: 0;
				font-size: 0.7rem;
				color: black;
			}
			input::placeholder {
				opacity: 1;
			}
		}
		
	</style>
	<div>
		<label for="input">${elem.attr("label")}</label>
		<input
			id="input" type="${elem.attr("type")}" name="${elem.name}" placeholder="${elem.attr("placeholder") || " " }"
			@input="${ evt => elem.value = evt.target.value}"
			@keypress="${ evt => elem.onKeyPress(evt)}"
		/>
	</div>
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
		this.#value = "";
		$.get(this.shadowRoot, "input").value = "";
	}
	onKeyPress(evt) {
		if (evt.code == "Enter" && this.#internal.form) {
			this.#internal.form.requestSubmit();
		}
	}
}
customElements.define("c-input", Input);
