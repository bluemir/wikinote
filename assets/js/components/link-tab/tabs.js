import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
			display: flex;
			gap: 0.3rem;
			border-bottom: 1px solid gray;
		}

		::slotted(.selected) {
			background: var(--blue-100);
		}
		::slotted(a) {
			padding: 0.3rem 0.5rem;
		}
	</style>
	<slot></slot>
`;

class CustomElement extends $.CustomElement {
	constructor() {
		super();
	}
	onConnected() {
		$.all(this, "a").
			filter((elem) => elem.hasAttribute("exact") ? elem.attr("href") == location.pathname : location.pathname.startsWith(elem.attr("href"))).
			forEach(elem => elem.classList.add("selected"));
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}
}
customElements.define("c-link-tabs", CustomElement);
