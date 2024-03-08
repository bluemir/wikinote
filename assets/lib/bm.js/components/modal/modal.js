import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (self) => html`
	<style>
		${css}

		:host {
			display: none;
			position: fixed;
			top: 0;
			left: 0;
			bottom: 0;
			right: 0;
		}
		:host(.show) {
			display: grid;
			background: #80808080;
		}
		#modal {
			margin: auto;
			background: #ffffff;
			min-width: 15rem;
		}
		#modal header {
			padding: 0.2rem 0.5rem;
			border-bottom: 1px solid gray;
		}
		#modal header ::slotted(*) {
			margin: 0;
			font-size: 1rem;
		}
		#modal section {
			min-height: 10rem;
			margin: 0.2rem 0.5rem;
		}
	</style>
	<section id="modal" @click="${evt => evt.stopPropagation()}">
		<header>
			<slot name="title">Modal</slot>
		</header>
		<section>
			<slot></slot>
		</section>
		<footer>
			<slot name="footer"></slot>
		</footer>
	</section>
`;

class Modal extends $.CustomElement {
	constructor() {
		super();
	}
	async onConnected() {
		this.on("click", evt => this.close())
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	open(){
		console.log("open")
		this.classList.add("show")
	}
	close() {
		this.classList.remove("show");
	}
}

customElements.define("c-modal", Modal);

