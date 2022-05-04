import * as $ from "../../lib/bm.js/bm.module.js";
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
			"logo menu"
			"nav  menu";

		/* FIXME */
		max-width: 1200px;
		margin: auto;
	}
	h1 {
		font-size: 2rem;
		margin: 0;
		padding: 0;

		grid-area: logo;
	}
	nav {
		grid-area: nav;
	}
	menu {
		grid-area: menu;

		margin: 0;
		padding: 0;
		justify-self: end;
		align-self: end;
	}
	form {
		grid-area: search;
		justify-self: end;
		padding: 0 0.7rem;
	}

	a {
		color:      var(--header-fg-color);
		text-decoration: none;
	}

	/* menu */
	menu ul {
		list-style: none;

		margin:  0px;
		padding: 0px;
	}

	menu a {
		color: var(--menu-font-color);
		white-space: nowrap;
	}
	menu ul li {
		padding: 0.3rem 0.7rem;
		display: inline-block;
		position: relative;
	}
	menu ul li:hover {
		background: var(--menu-hover-background-color);
	}

	menu ul.sub {
		display: none;
		position: absolute;
			right : 0rem;
		top: 100%;

		background: var(--sub-menu-background-color);
	}
	menu li:hover ul.sub {
		display: block;
	}

	menu ul.sub li {
		padding: 0.3rem 1.2rem;
		display: block;
	}
	menu ul.sub li:hover {
		background: var(--sub-menu-hover-background-color);
	}

	/* nav */
	nav {
		font-size: 0.8rem;
		margin-bottom: 0.3rem;
		color: var(--breadcrumbs-fg-color);
	}
	nav a {
		color: var(--breadcrumbs-fg-color);
	}
	nav a:hover {
		color: var(--breadcrumbs-hover-fg-color);
	}
</style>
<header>
	<h1>
		<a id="logo" href="/">Wikinote</a>
	</h1>
	<form action="/!/search" method="GET">
		<!-- TODO icon -->
		<input type="text" placeholder="Search" name="q"/><button type="submit">&#x1F50D;</button>
		<!--input type="submit" value="Search"/-->
	</form>
	<nav id="breadcrumbs">
		${app.breadcrumbs.map(item => html`/ <a href="${item.path}">${item.name}</a> `)}
	</nav>
	<menu>
		<!-- menu -->
		<ul>
			<li><a href="?edit">Edit</a></li>
			<li>
				<a href="">More</a>
				<ul class="sub">
					<li><a href="?delete">Delete</a></li>
					<hr>
					<!-- TODO make split -->
					<li><a href="/!/auth/logout">Logout</a></li>
					<li><a href="/!/login">Login</a></li>
					<li><a href="/!/auth/register">Sign Up</a></li>
					<li><a href="/!/user">Users</a></li>
				</ul>
			</li>
		</ul>
	</menu>
</header>
`;

class WikinoteHeader extends $.CustomElement {
	constructor() {
		super();

		this.on("connected", () => this.render())
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
