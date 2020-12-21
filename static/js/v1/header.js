import * as $ from "/!/static/lib/minilib.module.js";
import {html, render} from '/!/lib/lit-html/lit-html.js';

var tmpl = (app) => html`
<style>
	@import url("/!/static/css/color.css");
	:host {
		display: block;

		background: var(--header-background-color);
		color: var(--header-font-color);
	}
	header {
		display: grid;
		grid-template-columns: 1fr auto;
		grid-template-areas:
			"logo search"
			"logo menu"
			"nav  menu";
	}
	h1 {
		font-size: 1.3rem;
		margin: 0;
		padding: 0;

		grid-area: logo;
	}
	nav {
		grid-area: nav;
	}
	menu {
		grid-area: menu;
	}
	form {
		grid-area: search;
		justify-self: end;
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
		{{ range .breadcrumb }}
		/ <a href="{{.Path}}">{{ .Name }}</a>
		{{ end }}
	</nav>
	<menu>
		<!-- menu -->
		{{ template "menu" . }}
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
}
customElements.define("wikinote-header", WikinoteHeader);
