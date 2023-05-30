import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		header {
			display: grid;
			justify-content: space-between;
			grid-template-columns: auto 1fr auto;
		}

		header p {
			margin: 0px;
		}

		:host([state=open]) section {
			display: block;
		}
		:host section {
			display: none;
			padding-left: 1rem;
			margin-top: 0.5rem;
		}
	</style>
	<header @click="${evt => elem.toggle(evt)}">
		<c-icon kind="${elem.state == "open" ? "expand_less": "expand_more"}"></c-icon>
		<p>${elem.attr("title") || "extends"}</p>
		<slot name="more"></slot>
	</header>
	<section>
		<slot></slot>
	</section>
`;

class Collapse extends $.CustomElement {
	constructor() {
		super();
	}
	static get observedAttributes() {
		return ["title", "state"];
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	onAttributeChanged(name, old, v) {
		this.render();
	}
	toggle() {
		this.state = this.state == "open" ? "close": "open"
	}
	// attribute
	get state() {
		return this.attr("state");
	}
	set state(v) {
		this.attr("state", v)
	}
}
customElements.define("c-collapse", Collapse);

