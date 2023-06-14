import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");
		@import url("/!/static/css/markdown.css");

		:host {
			display: block;
			min-height: 30rem;
		}
	</style>
	<div class="placeholder">
	</div>
`;

class WikinoteViewerMarkdown extends $.CustomElement {
	constructor() {
		super();

		this.on("connected", () => this.onConnected())
	}

	onConnected() {
		this.render();
		// slot? or dom migrate?,
		//$.get(this.shadow, ".placeholder").innerHTML = this.innerHTML;
		this.shadow.append(...this.childNodes);
		// https://stackoverflow.com/questions/61626493/slotted-css-selector-for-nested-children-in-shadowdom-slot
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute

	// event listener
}
customElements.define("wikinote-viewer-markdown", WikinoteViewerMarkdown);
