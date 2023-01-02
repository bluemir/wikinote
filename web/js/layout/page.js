import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		:host {
			position: fixed;
			display: grid;
			grid-template-rows: auto 1fr;
			width: 100%;
			height: 100%;
		}
		wikinote-header {

		}
		wikinote-header::part(wrapper), main {
			padding: 0 2rem;
			max-width: 1200px;
			margin: 0 auto;
		}
		section {
			overflow-y: scroll;

			background: var(--contents-padding-bg-color);
		}

		main {
			width: 100%;
			padding: 2rem;
			background: var(--contents-bg-color);
		}
	</style>
	<wikinote-header></wikinote-header>
	<section>
		<main>
			<slot></slot>
		</main>
	</section>
`;

class WikinotePage extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-page", WikinotePage);
