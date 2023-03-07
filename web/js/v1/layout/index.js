import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

import "./header.js";
import "./page.js";

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		:host {
			display: block;
			background: var(--contents-padding-bg-color);

			margin:  0px;
			padding: 0px;

			overflow: hidden;
			top: 0;
			left: 0;
			right: 0;
			bottom: 0;
			position: fixed;

			/* grid */
			display: grid;
			grid-template-areas: "header" "page";
			grid-template-rows: auto 1fr;
		}
		wikinote-header {
			grid-area: "header";
		}
		wikinote-page {
			grid-area: "page";
		}
	</style>
	<wikinote-header></wikinote-header>
	<wikinote-page>
		<slot></slot>
	</wikinote-page>
`;

class WikinoteLayout extends $.CustomElement {
	constructor() {
		super();

		this.render();
		this.on("connected", () => this.render())
	}

	render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("wikinote-layout", WikinoteLayout);
