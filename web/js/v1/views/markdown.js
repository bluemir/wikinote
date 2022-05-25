import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		:host {
			display: block;
			min-height: 30rem;
		}
		a {
			color: var(--link-fg-color);
			text-decoration: none;
		}
		a:hover {
			text-decoration: underline;
		}
		h1 {
			font-size: 2rem;
			text-shadow: 1px 1px 5px var(--gray);
		}
		h2 {
			font-size: 1.7rem;
			border-bottom: 1px solid var(--gray);
		}
		h3 {
			font-size: 1.5rem;
		}
		h4 {
			font-size: 1.3rem;
		}
		h5 {
			font-size: 1.2rem;
		}
		table {
			border-collapse: collapse;
		}
		table th, table td {
			border: 1px solid var(--gray);
		}
	</style>
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
