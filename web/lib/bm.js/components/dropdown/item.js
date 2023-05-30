import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		:host {
			display: block;
			padding: 0.3rem 1rem;
		}
		:host(:hover) {
			background: #343434;
		}

		::slotted(*) {
			display: block;
			text-decoration: none;
			color: white;
			white-space: nowrap;
		}
	</style>
	<slot></slot>
`;

class DropdownItem extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-dropdown-item", DropdownItem);
