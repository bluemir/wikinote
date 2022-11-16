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

		#logo a {
			color:      var(--header-fg-color);
			text-decoration: none;
		}

	</style>
	<header>
		<section id="logo">
			<a href="/">Wikinote</a>
		</section>
		<section class="search">
		</section>
		<nav>
			${app.breadcrumbs.map(item => html`/ <a href="${item.path}">${item.name}</a> `)}
		</nav>
		<menu>
			<ul>
				<li><a href="?edit">Edit</a></li>
				<li>
					<a href="">More</a>
					<ul class="sub">
						<li><a href="?delete">Delete</a></li>
						<hr>
						<!-- TODO make split -->
						<li><a href="/!/auth/login">Login</a></li>
						<li><a href="/!/auth/profile">Profile</a></li>
						<li><a href="/!/auth/register">Sign Up</a></li>
					</ul>
				</li>
			</ul>
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
