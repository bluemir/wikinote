import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

var tmpl = (app) => html`
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

class WikinoteHeaderLogo extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadowRoot);
	}
	// attribute
}
customElements.define("wikinote-header-logo", WikinoteHeaderLogo);
