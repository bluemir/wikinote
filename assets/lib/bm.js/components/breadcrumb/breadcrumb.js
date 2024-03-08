import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<style>
		nav {
			font-size: 0.8rem;
			color: var(--breadcrumbs-fg-color);
			padding: 0px;
		}
		nav a {
			text-decoration: none;
			color: var(--breadcrumbs-fg-color);
		}
		nav a:hover {
			text-decoration: underline;
			color: var(--breadcrumbs-hover-fg-color);
		}
	</style>
	<nav>
		${elem.breadcrumbs.map(item => html`/ <a href="${item.path}">${item.name}</a> `)}
	</nav>
`;

class Breadcrumbs extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	// attribute
	get breadcrumbs() {
		const arr = location.pathname.split("/").filter(e => e.length);

		return arr.map((item, index) => {
			return {
				name: item,
				path: "/" + arr.slice(0, index+1).join("/"),
			}
		});
	}
}
customElements.define("c-breadcrumbs", Breadcrumbs);

