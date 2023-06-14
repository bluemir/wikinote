import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (app) => html`
	<style>
		${css}

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
