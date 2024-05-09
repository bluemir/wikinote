import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

var tmpl = (app) => html`
	<style>
		${css}
		
		:host {
			display: block;
		}
	</style>
	<slot></slot>
`;

class WikinotePluginFooter extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-plugin-footer", WikinotePluginFooter);
