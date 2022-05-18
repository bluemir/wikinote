import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		:host {
			display: block;
			min-height: 30rem;
		}
	</style>
	<img src=${app.attr("url")} />
`;

class WikinoteViewImage extends $.CustomElement {
	constructor() {
		super();

		this.on("connected", () => this.onConnected())
	}

	onConnected() {
		this.render();
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute

	// event listener
}
customElements.define("wikinote-view-image", WikinoteViewImage);
