import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<style>
		@import url("/-/static/css/system/grid.css");
		
		:host {
			display: block;
		}
	</style>
	<h1>Profile</h1>
	<section grid>
		<label col="2">Name</label>
		<article col="10">${elem.user?.name}</article>
	</section>
	<section grid>
		<label col="2">Groups</label>
		<article col="10">${elem.user?.groups?.join(", ")}</article>
	</section>
	<!-- TODO labels -->
`;

class CustomElement extends $.CustomElement {
	constructor() {
		super();

	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute

	onConnected() {
		this.user = $.get(this, "c-data").json;
		this.render();
	}
}
customElements.define("wikinote-profile", CustomElement);
