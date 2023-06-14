import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (app) => html`
	<style>
		${css}

		:host {
			display: inline-block;
			padding: 0.3rem 0.8rem;
		}
		:host(:hover) {
			background: var(--bg-color, #343434);
		}

		::slotted(*) {
			display: block;
			text-decoration: none;
			color: var(--fg-color, white);
			white-space: nowrap;
		}
	</style>
	<slot></slot>
`;

class Button extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-button", Button);
