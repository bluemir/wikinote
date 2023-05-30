import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

/*

<c-dropdown title="title">
	<c-dropdown-item>a</c-dropdown-item>
	<c-dropdown-item>b</c-dropdown-item>
</c-dropdown>

*/

var tmpl = (app) => html`
	<style>
		:host {
			display: inline-block;
			position: relative;
			padding: 0.3rem 0.8rem;
		}
		a {
			color: var(--header-fg-color);
			text-decoration: none;
		}
		menu {
			display: none;
			position: absolute;
			right: 0;
			top: 100%;

			background: gray;
			padding: 0rem;
			margin: 0rem;
		}
		::slotted(dropdown-item) {
			display: block;
		}
		:host(:hover) menu {
			display: block;
		}
	</style>
	<a href="">${app.attr("title")}</a>
	<menu>
		<slot></slot>
	</menu>
`;

class Dropdown extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-dropdown", Dropdown);
