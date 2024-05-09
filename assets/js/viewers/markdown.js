import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

var tmpl = (app) => html`
	<style>
		${css}

		:host {
			display: block;
			min-height: 35rem;
		}
	</style>
	<slot></slot>
`;

class WikinoteViewerMarkdown extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-viewer-markdown", WikinoteViewerMarkdown);
