import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
<style>
	@import url("/!/static/css/color.css");
</style>
`;

class WikinoteSearch extends $.CustomElement {
	constructor() {
		super();

		this.on("connected", () => this.render())
	}
	onConnected() {
		this.shadow.append(...this.childNodes);
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("wikinote-search", WikinoteSearch);
