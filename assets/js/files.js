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
	<input type="file" />
	<ul>
		${app.files?.map(file => html`
			<li><a href="${file.path}">${file.name}</a></li>
		`)}
	</ul>
`;

class CustomElement extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
	get files() {
		return $.get(this, "c-data").json;
	}
}
customElements.define("wikinote-files", CustomElement);
