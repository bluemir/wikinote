import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		:host {
			display: grid;
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
