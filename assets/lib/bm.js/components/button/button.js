import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
			display: inline-block;

			color: var(--button-fg-color, black);
			background: var(--button-bg-color, white);
			border: 1px solid var(--button-border-color, #343434);
		}
		:host(:hover) {
			color: var(--button-hover-fg-color, white);
			background: var(--button-hover-bg-color, #343434);
			border: 1px solid var(--button--hover-border-color, #343434);
		}
		::slotted(a) {
			color: inherit;
			text-decoration: none;
			white-space: nowrap;
		}
		::slotted(a:hover) {
		}
		::slotted(button) {
			background: none;
			font: inherit;
			border: 0;
			padding: 0;
		}
		:host(:hover) ::slotted(button) {
			color: var(--button-hover-fg-color, white);
		}

		/* ':host(:has(a, button))' is not available for chrome & firefox. it seems bug.
		 * see also: https://github.com/w3c/webcomponents-cg/issues/5#issuecomment-1220786480
		 */
		:host([pure]) {
			padding: 0.3rem 0.8rem;
		}
		::slotted(a), ::slotted(button) {
			display: inline-block;
			padding: 0.3rem 0.8rem;
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
	onConnected() {
		if(!$.get(this, "a, button")) {
			this.setAttribute("pure", "");
		}
	}
}
customElements.define("c-button", Button);
