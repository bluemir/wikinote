import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		:host {
			display: grid;
			align-items: center;
			justify-items: center;

			padding: 0.8rem;
			border-radius: 0.5rem;
			background: #9e9e9e;
			color: #ffffff;
			font-weight: bold;
			${elem.width};
			${elem.height};
		}
	</style>
	${elem.attr("text") || "placeholder"}
`;

class Placeholder extends $.CustomElement {
	constructor() {
		super();
	}
	static get observedAttributes() {
		return ["text", "width", "height"];
	}
	onAttributeChanged() {
		this.render();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	get width() {
		let w = this.attr("width");
		return w ? `width: ${w}` : ``;
	}
	get height() {
		let h = this.attr("height");
		return h ? `height: ${h}`: ``;
	}
}
customElements.define("c-placeholder", Placeholder);
