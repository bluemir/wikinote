import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

// TODO
// - animation
// - proxy attribute of input
// - type of input(outline, underline, ...)
class Input extends $.CustomElement {
	template() {
		return html`
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
				<label for="input">${this.attr("label")}</label>
				<input
					id="input" type="${this.attr("type")}" name="${this.name}" placeholder="${this.attr("placeholder") || " " }"
					@input="${ evt => this.value = evt.target.value}"
					@keypress="${ evt => this.onKeyPress(evt)}"
				/>
			</div>
		`;
	}

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
		render(this.template(), this.shadowRoot);
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
		$.get(this.shadow, "input").value = "";
	}
	onKeyPress(evt) {
		if (evt.code == "Enter" && this.#internal.form) {
			this.#internal.form.requestSubmit();
		}
	}
}
customElements.define("c-input", Input);
