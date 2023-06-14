import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		:host {
			display: block;
		}
	</style>
	<h1>Messages</h1>
	${app.messages.map(msg => html`
		<p>
			${msg.At} - ${msg.Detail.Text}
		</p>
	`)}
`;

class WikinoteMessages extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		this.messages = $.get(this, "c-data").json;

		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-messages", WikinoteMessages);
