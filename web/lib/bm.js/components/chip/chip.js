import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
			display: flex;
			padding: 0.5rem;
			border: 1px solid gray;
			border-radius: 0.5rem;
		}
	</style>
	<slot></slot>
`;

class Chip extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-chip", Chip);
