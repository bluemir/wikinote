import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (app) => html`
	<style>
		${css}

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
