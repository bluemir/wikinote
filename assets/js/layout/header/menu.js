import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {me, can} from "api.js";

var tmpl = (elem) => html`
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
			${elem.me?html`
				<c-dropdown-item>
					<a href="/-/auth/login?exclude=${elem.me.name}">Logout</a>
				</c-dropdown-item>
			`:html`
				<c-dropdown-item>
					<a href="/-/auth/login">Login</a>
				</c-dropdown-item>
				<c-dropdown-item>
					<a href="/-/auth/register">Register</a>
				</c-dropdown-item>
			`}
			<c-dropdown-item>
				<a href="/-/auth/profile">Profile</a>
			</c-dropdown-item>
			<c-dropdown-item>
				<a href="/-/messages">Messages</a>
			</c-dropdown-item>
			${elem.canAccessAdmin?html`<c-dropdown-item><a href="/-/admin">Admin</a></c-dropdown-item>`:""}
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
	async onConnected(){
		this.me = await me();
		this.canAccessAdmin = await can("read", "admin");

		this.render();
	}
	
}
customElements.define("wikinote-header-menu", WikinoteHeaderMenu);
