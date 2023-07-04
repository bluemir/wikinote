import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined" rel="stylesheet" />
	<style>
		span.material-symbols-outlined {
			${elem.size}
			cursor: default;
		}
	</style>

	<span class="material-symbols-outlined">${elem.attr("kind")}</span>
`;

class Icon extends $.CustomElement {
	constructor() {
		super();
	}
	onConnected() {
		// add stylesheet
		if ($.get(document, "head link#icons")) {
			// skip
			return
		}
		$.get(document, "head").appendChild($.create("link", {
			href: "https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined",
			rel: "stylesheet",
			id: "icons",
		}));
	}
	static get observedAttributes() {
		return ["kind", "size"];
	}
	onAttributeChanged(name, old, v) {
		this.render();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
	get size() {
		let n = this.attr("size");
		return n ? `font-size: ${n};` : ""
	}
}
customElements.define("c-icon", Icon);

