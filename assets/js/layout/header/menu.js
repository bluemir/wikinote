import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import * as auth from "../../auth.js";

var tmpl = (app) => html`
	<style>
		@import url("/-/static/css/color.css");

		menu {
			padding: 0px;
			margin: 0px;
		}
	</style>
	<menu>
		<c-button>
			<a href="?edit">Edit</a>
		</c-button>
		<c-dropdown title="More">
			<c-dropdown-item>
				<a href="?files">Files</a>
			</c-dropdown-item>
			<hr />
			<c-dropdown-item>
				<a href="/-/auth/login">Login</a>
			</c-dropdown-item>
			<c-dropdown-item>
				<a href="/-/auth/profile">Profile</a>
			</c-dropdown-item>
			<c-dropdown-item>
				<a href="/-/messages">Messages</a>
			</c-dropdown-item>
			<c-dropdown-item>
				<a href="/-/auth/login">Sign Up</a>
			</c-dropdown-item>
			${auth.can("read", "admin")?html`<c-dropdown-item><a href="/-/admin">Admin</a></c-dropdown-item>`:""}
		</c-dropdown>
	</menu>
`;

class WikinoteHeaderMenu extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-header-menu", WikinoteHeaderMenu);
