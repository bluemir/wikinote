import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");
	</style>
	<h1>${ app.attr("user") }</h1>
	<h2>Roles</h2>
	<ul>
		${ app.roles.map(role => html`
			<li>${role}</li>
		`)}
	</ul>
	<button @click=${evt => app.logout(evt)}>Logout</button>
`;

class WikinoteProfile extends $.CustomElement {
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
	get username() {
		return this.attr("user")
	}
	get roles() {
		return this.attr("roles").split(",")
	}

	// event listener
	logout(evt) {

	}
}
customElements.define("wikinote-profile", WikinoteProfile);
