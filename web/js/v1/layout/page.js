import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

import * as api from "v1/api.js";

var tmpl = (app) => html`
<style>
	@import url("/!/static/css/color.css");
	:host {
		display: block;
		overflow-y: scroll;
	}
	section {
		max-width: 1200px;
		margin: auto;
		margin-top: 0;
		padding: 2rem;

		background: var(--contents-bg-color);
	}
</style>
<section>
	<slot></slot>
</section>
`;
// TODO scroll bar custom

class WikinotePage extends $.CustomElement {
	constructor() {
		super();

		this.on("connected", () => this.render())
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-page", WikinotePage);
