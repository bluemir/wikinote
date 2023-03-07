import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");
	</style>
	<form method="post">
		<dl>
			<dt>
				<label for="username">Username</label>
			</dt>
			<dd>
				<input name="username" placeholder="username"/>
			</dd>
			<dt>
				<label for="email">Email</label>
			</dt>
			<dd>
				<input name="email" placeholder="email" type="email" />
			</dd>
			<dt>
				<label for="password">Password</label>
			</dt>
			<dd>
				<input name="password" type="password" placeholder="password"/>
			</dd>
			<dt>
				<label for="confirm">Password Confirm</label>
			</dt>
			<dd>
				<input name="confirm" type="password" placeholder="password confirm"/>
			</dd>

		</dl>
		<button type="submit">Submit</button>
	</form>
`;

class WikinoteRegister extends $.CustomElement {
	constructor() {
		super();
	}

	onConnected() {
		this.render();
	}

	async render() {
		render(tmpl(this), this.shadow);
	}

	// attribute
	get path() {
		return this.attr("path");
	}

	// event listener
}
customElements.define("wikinote-register", WikinoteRegister);
