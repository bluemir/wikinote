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
					}
					
				}
				::slotted(a) {
					color: white;
					display: block;
					text-decoration: none;
					padding: 0.5rem;
				}
			</style>
			<section dropdown>
				<a href="#">${title}</a>
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
