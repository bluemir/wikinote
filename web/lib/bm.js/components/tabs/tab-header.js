import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		:host {
			background: lightgray;
			padding: 0.3rem 0.5rem;
		}
		:host(.selected) {
			background: gray;
		}
	</style>
	<slot></slot>
`;

class TabHeader extends $.CustomElement {
	constructor() {
		super();
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-tab-header", TabHeader);
