import * as $ from "bm.js/bm.module.js";
//import {html, render} from '/lib/lit-html/lit-html.js';
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<style>
		@import "/static/css/root.css";

		:host {
			display: inline-block;
		}
		::slotted(*) {
		}
		section.container {
			display: none;
			position: absolute;

			background: white;
			padding: 1rem;
		}
		a:hover {
			section {
				display: block;
			}
		}
	</style>
	<a>
		<slot name="trigger" @mouseover="${evt => elem.onHover(evt)}"></slot>
		<section class="container" shadow="2">
			<slot></slot>
		</section>
	</a>
`;

class Popover extends $.CustomElement {
	static get observedAttributes() {
		return [];
	}
	constructor() {
		super();
	}

	onAttributeChanged(name, oValue, nValue) {
	}

	onConnected() {
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}
	onHover(evt) {
		console.log(evt);
	}
}
customElements.define("c-popover", Popover);
