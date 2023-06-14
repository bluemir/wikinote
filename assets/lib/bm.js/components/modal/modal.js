import * as $ from "bm.js/bm.module.js";

import {html, render} from 'lit-html';

var tmpl = (self) => html`
	<style>
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
		#modal section {
			min-height: 10rem;
		}
	</style>
	<section id="modal" @click="${evt => evt.stopPropagation()}">
		<header><slot name="title"></slot></header>
		<section>
			<slot name="main"></slot>
		</section>
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

