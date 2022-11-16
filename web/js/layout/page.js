import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		:host {
			display: block;
		}
	</style>
	<wikinote-header></wikinote-header>
	<main>
		<slot></slot>
	</main>
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
