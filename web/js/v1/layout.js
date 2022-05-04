import * as $ from "../../lib/bm.js/bm.module.js";
import {html, render} from 'lit-html';

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
		.page {
			overflow-y: scroll;

			grid-area: "page";
		}
		article {
			max-width: 1200px;
			margin: auto;
			margin-top: 0;
			padding: 2rem;

			background: var(--contents-bg-color);
		}
	</style>
	<wikinote-header></wikinote-header>
	<div class="page">
		<article>
			<slot></slot>
		</article>
	</div>
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
