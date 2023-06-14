import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (app) => html`
	<style>
		${css}

		:host {
			display: grid;

			padding: 0.5rem;

			box-shadow: 2px 2px 4px 2px rgba(0,0,0,0.5);
		}

		header {
			border-bottom: 1px solid rgba(0,0,0,0.5);
				margin-bottom: 0.5rem;
		}
	</style>
	<section>
		<header>
			<slot name="header"></slot>
		</header>
		<div>
			<slot name="body"></slot>
		</div>
	</section>
`;

class Card extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
}
customElements.define("c-card", Card);
