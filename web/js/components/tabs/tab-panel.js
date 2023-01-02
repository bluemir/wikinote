import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/static/css/common.css");

		:host {
		}
	</style>
	<slot></slot>
`;

class TabPanel extends $.CustomElement {
	constructor() {
		super();
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-tab-panel", TabPanel);
