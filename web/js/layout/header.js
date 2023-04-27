import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import * as auth from "../auth.js";

var tmpl = (app) => html`
	<style>
		@import url("/-/static/css/color.css");
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

		wikinote-header-logo {
			grid-area: logo;
		}
		wikinote-header-search {
			grid-area: search;
		}
		/* nav */
		c-breadcrumbs {
			grid-area: nav;
			align-self: end;
			margin-bottom: 0.3rem;
		}
		wikinote-header-menu {
			grid-area: menu;
			align-self: end;
			justify-self: end;
		}
	</style>
	<header part="wrapper">
		<wikinote-header-logo></wikinote-header-logo>
		<wikinote-header-search></wikinote-header-search>
		${ app.isSpecialPath ? "": html`<c-breadcrumbs></c-breadcrumbs>` }
		<wikinote-header-menu></wikinote-header-menu>
	</header>
`;

//

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
			case "-":
			case "~":
				return true;
		}
		return false;
	}
}
customElements.define("wikinote-header", WikinoteHeader);
