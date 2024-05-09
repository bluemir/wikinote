import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {shortcut} from "shortcut.js";
import {css} from "common.js";

var tmpl = (app) => html`
	<style>
		${css}

		:host {
			overflow-y: scroll;
		}
		wikinote-header::part(wrapper), main {
			padding: 0 2rem;
			max-width: 1200px;
			margin: 0 auto;
		}

		wikinote-header {
			position: sticky;
			top: 0;
			width: 100%;
		}

		main {
			padding: 2rem;
			background: var(--contents-bg-color);
		}
	</style>
	<wikinote-header></wikinote-header>
	<main>
		<slot></slot>
	</main>

`;

class CustomElement extends $.CustomElement {
	constructor() {
		super();
	}
	onConnected() {
		shortcut.add(`ctrl+shift+e`, evt => {
			// show editor
			console.log("editor");

			location.assign(location.pathname+"?edit");
		})
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-page", CustomElement);
