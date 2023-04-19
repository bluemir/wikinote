import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/-/static/css/color.css");
		a {
			color: var(--header-fg-color);
			text-decoration: none;
			font-size: 2rem;
			font-weight: bold;
		}
	</style>
	<a href="/">Wikinote</a>
`;

class WikinoteLogo extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-logo", WikinoteLogo);
