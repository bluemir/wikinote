import * as $ from "/!/static/lib/minilib.module.js";
import {html, render} from '/!/lib/lit-html/lit-html.js';

var tmpl = (app) => html`
<style>
	:host {
	}
</style>
<wikinote-header></wikinote-header>
<slot>
</slot>
`;

class WikinoteLayout extends $.CustomElement {
	constructor() {
		super();

		this.on("connected", () => this.render())
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("wikinote-layout", WikinoteLayout);
