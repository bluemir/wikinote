import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");
		*:not(:defined) {
			display:none;
		}

		:host {
			display: block;
			background: var(--header-bg-color);
			color:      var(--header-fg-color);
		}

		header {
			display: grid;
			grid-template-columns: 1fr auto;
			grid-template-areas:
				"logo search"
				"nav  menu";
		}

		#logo {
			grid-area: logo;
		}

		#logo a {
			color:      var(--header-fg-color);
			text-decoration: none;
			font-size: 2rem;
			font-weight: bold;
		}
		.search {
			grid-area: search;
		}
		/* nav */
		c-breadcrumbs {
			grid-area: nav;
			align-self: end;
			margin-bottom: 0.3rem;
		}

		menu {
			grid-area: menu;
			align-self: end;
			justify-self: end;
			padding: 0px;
			margin: 0px;
		}

	</style>
	<header part="wrapper">
		<section id="logo">
			<a href="/">Wikinote</a>
		</section>
		<section class="search">
			<form action="/!/search">
				<input name="q"/>
				<button><c-icon kind="search"></c-icon></button>
			</form>
		</section>
		${ app.isSpecialPath ? "": html`<c-breadcrumbs></c-breadcrumbs>` }
		<menu>
			<c-button>
				<a href="?edit">Edit</a>
			</c-button>
			<c-dropdown title="More">
				<c-dropdown-item>
					<a href="/!/auth/login">Login</a>
				</c-dropdown-item>
				<c-dropdown-item>
					<a href="/!/auth/profile">Profile</a>
				</c-dropdown-item>
				<c-dropdown-item>
					<a href="/!/auth/login">Sign Up</a>
				</c-dropdown-item>
			</c-dropdown>
		</menu>
	</header>
`;

class WikinoteHeader extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
	get isSpecialPath() {
		const arr = location.pathname.split("/").filter(e => e.length);
		switch (arr[0]) {
			case "!":
			case "%21":
				return true;
		}
		return false;
	}
}
customElements.define("wikinote-header", WikinoteHeader);
