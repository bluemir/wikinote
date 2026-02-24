import * as $ from "bm.js/bm.module.js";
//import {html, render} from '/lib/lit-html/lit-html.js';
import {html, render} from 'lit-html';

var tmpl = (self) => html`
	<style>
		@import "/static/css/root.css";

		:host {
		}
	</style>
	<a href="#" @click="${evt => self.goBack(evt)}">Go Back</a>
`;

class GoBack extends $.CustomElement {
	constructor() {
		super();
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}

	async goBack(evt) {
		evt.preventDefault();

		history.back();
	}
}
customElements.define("c-go-back", GoBack);
