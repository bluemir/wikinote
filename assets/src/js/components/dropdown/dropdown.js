import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";

class CustomElement extends $.CustomElement {
	constructor() {
		super();
	}

	template(title){
		return html`
			<style>
				${css}

				:host {
					display: block;
				}

				[dropdown] {
					position: relative;

					& > a {
						color: white;
						text-decoration: none;
						display: block;
						padding: 0.5rem;
					}
					& > section {
						display: none;
						
					}
					&:hover > section {
						display: block;
						position: absolute;
						background: gray;
						right: 0;
					}
				}
				/* menu item */
				/* TODO color */
				::slotted(a) {
					color: white;
					display: block;
					text-decoration: none;
					padding: 0.5rem;
				}
				::slotted(a:hover) {
					background: lightgray;
				}
			</style>
			<section dropdown>
				<a href="#" part="anchor">${title}</a>
				<section>
					<slot></slot>
				</section>
			</section>
		`;
	}

	async render() {
		render(this.template(this.attr("name")), this.shadowRoot);
	}
}
customElements.define("c-dropdown", CustomElement);
