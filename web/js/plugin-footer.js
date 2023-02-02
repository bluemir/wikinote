import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
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
