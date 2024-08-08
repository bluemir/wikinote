import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<style>
		@import "/static/css/root.css";

		:host {
		}
 	</style>
	<input type="range" /><input type="number" />
`;

class Slider extends $.CustomElement {
	constructor() {
		super();
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}
}
customElements.define("c-slider", Slider);
