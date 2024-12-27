import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

var tmpl = (app) => html`
	<style>
		${css}

		:host {

		}
	</style>
	<form action="/-/search">
		<input name="q" value="${app.query}" />
		<button><c-icon kind="search" size="1rem"></c-icon></button>
	</form>
`;

class CustomElement extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadowRoot);
	}
	// attribute
	get query() {
		//if (location.pathname != "/-/search") {
		//	return ""
		//}

		let params = new URLSearchParams(document.location.search);

		return params.get("q");
	}
}
customElements.define("wikinote-header-search", CustomElement);
