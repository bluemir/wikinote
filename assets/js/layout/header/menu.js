import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {me, can} from "api.js";
import {css} from "common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
			display: flex;
		}
		menu {
			padding: 0px;
			margin: 0px;
			display: block;
		}
		a {
			color: white;
			padding: 0.5rem 1rem;
			text-decoration: none;
			display: block;
		}
		a:hover {
			background-color: var(--gray-700);
		}

		menu[dropdown] {
			[role=dropdown] {
				display: none;
				
				position: absolute;
				background: var(--gray-800);
			}
		}
		menu[dropdown]:has(:hover) {
			[role=dropdown] {
				display: block;
			}
		}
	</style>
	<a href="?edit">Edit</a>
	<menu dropdown>
		<a href="#" role="trigger">More</a>
		<section role="dropdown">
			<a href="?files">Files</a>
			<hr />
			${elem.me?html`
				<a href="/-/auth/profile">Profile</a>
				<a href="/-/auth/login?exclude=${elem.me.name}">Logout</a>
				<a href="/-/messages">Messages</a>
			`:html`
				<a href="/-/auth/login">Login</a>
				<a href="/-/auth/register">Register</a>
			`}
			${elem.canAccessAdmin?html`<a href="/-/admin">Admin</a>`:""}
		</section>
	</menu>
`;

class CustomElement extends $.CustomElement {
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
customElements.define("wikinote-header-menu", CustomElement);
