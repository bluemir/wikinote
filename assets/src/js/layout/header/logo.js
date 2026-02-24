import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

class CustomElement extends $.CustomElement {
	template() {
		return html`
			<style>
				${css}

				a {
					color: var(--header-fg-color);
					text-decoration: none;
					font-size: 2rem;
					font-weight: bold;
				}
			</style>
			<a href="/">Wikinote</a>
		`;
	}
	constructor() {
		super();
	}
	async render() {
		render(this.template(), this.shadowRoot);
	}
	// attribute
}
customElements.define("wikinote-header-logo", CustomElement);
