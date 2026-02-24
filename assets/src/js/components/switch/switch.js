import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<style>
		@import "/static/css/root.css";

		:host {
			display: inline-block;
			margin: 0.3rem 0.5rem;
			height: 1rem;
			width: 2rem;
			vertical-align: middle;
		}
		#bar {
			height: 100%;
			width: 100%;
			background: var(--gray-200);
			border: 1px solid var(--gray-400);
			border-radius: 100vh;

			& > #thumb {
				aspect-ratio : 1 / 1;
				height: 100%;
				background: var(--blue-500);
				border-radius: 100vh;
			}
		}
	</style>
	<div id="bar">
		<div id="thumb"></div>
	</div>
`;

class Switch extends $.CustomElement {
	static get observedAttributes() {
		return [];
	}
	constructor() {
		super();
	}

	onAttributeChanged(name, oValue, nValue) {
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}
}
customElements.define("c-switch", Switch);
