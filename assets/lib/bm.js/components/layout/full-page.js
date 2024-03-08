import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
			position: fixed;
			width: 100%;
			height: 100%;

			display: grid;
			grid-template-columns: auto 1fr;
			grid-template-rows: auto 1fr auto;
		}
		header {
			background: var(--header-bg-color, black);
			padding: 0 1rem;
			grid-column: span 2;
		}
		header ::slotted(*) {
			--fg-color: var(--header-fg-color, white);
		}
		aside {
			border: 1px solid var(--border-color, gray)
		}
		main {
			padding: 1rem;
			overflow-y: scroll;
		}
		footer {
			grid-column: span 2;
		}
	</style>
	<header>
		<slot name="header"></slot>
	</header>
	<aside>
		<slot name="menu"></slot>
	</aside>
	<main>
		${elem.hasAttribute("footer") ? html`
			<slot name="main"></slot>
		` : html`
			<slot></slot>
		`}
	</main>
	<footer>
		<slot name="footer"></slot>
	</footer>
`;

class FullPage extends $.CustomElement {
	constructor() {
		super();
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-layout-full-page", FullPage);
