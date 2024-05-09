import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

var tmpl = (app) => html`
	<style>
		${css}

		:host {
			display: block;
			min-height: 30rem;
		}
		img {
			max-width: 100%;
		}
	</style>
	<video src=${app.attr("url")} controls ></video>
`;

class WikinoteViewerVideo extends $.CustomElement {
	constructor() {
		super();
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute

	// event listener
}
customElements.define("wikinote-viewer-video", WikinoteViewerVideo);
