import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
			display: block;
			font-size: 1rem;
			padding: 0.8rem;
			border-radius: 0.5rem;
			background: #2fcc66;
			color: #ffffff;
			font-weight: bold;
		}
	</style>
	${elem.attr("text")}
`;

class Alert extends $.CustomElement {
	constructor() {
		super();
	}
	static get observedAttributes() {
		return ["text"];
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-alert", Alert);
