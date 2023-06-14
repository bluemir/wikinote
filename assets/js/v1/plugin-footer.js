import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		:host {
			display: block;
		}
	</style>
`;

class WikinotePluginFooter extends $.CustomElement {
	constructor() {
		super();

		this.on("connected", () => this.onConnected())
	}

	onConnected() {
		this.render();
		this.shadow.append(...this.childNodes);
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute

	// event listener
}
customElements.define("wikinote-plugin-footer", WikinotePluginFooter);
