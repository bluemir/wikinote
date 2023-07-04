import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/-/static/css/color.css");

	</style>
	<form action="/-/search">
		<input name="q" value="${app.query}" />
		<button><c-icon kind="search" size="1rem"></c-icon></button>
	</form>
`;

class WikinoteHeaderSearch extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
	get query() {
		if (location.pathname != "/-/search") {
			return ""
		}

		let params = new URLSearchParams(document.location.search);

		return params.get("q");
	}
}
customElements.define("wikinote-header-search", WikinoteHeaderSearch);
