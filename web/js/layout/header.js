import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

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

		#logo a {
			color:      var(--header-fg-color);
			text-decoration: none;
			font-size: 2rem;
			font-weight: bold;
		}
		/* nav */
		c-breadcrumbs {
			align-self: end;
			margin-bottom: 0.3rem;
		}

		*:not(:defined) {
			display:none;
		}

		menu {
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
			<form>
				<input />
				<button><c-icon kind="search"></c-icon></button>
			</form>
		</section>
		<c-breadcrumbs></c-breadcrumbs>
		<menu>
			<c-button>
				<a href="?edit">Edit</a>
			</c-button>
			<c-dropdown title="More">
				<c-dropdown-item>
					<a href="/!/auth/login">Login</a>
				</c-dropdown-item>
				<c-dropdown-item>
					<a href="/!/auth/login">Profile</a>
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
	get breadcrumbs() {
		const arr = location.pathname.split("/").filter(e => e.length);
		if ( arr[0] == "!") {
			return [];
		}

		return arr.map((item, index) => {
			return {
				name: item,
				path: "/" + arr.slice(0, index+1).join("/"),
			}
		});
	}
}
customElements.define("wikinote-header", WikinoteHeader);
