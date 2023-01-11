import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

class Data extends $.CustomElement {
	constructor() {
		super({enableShadow: false});
	}
	// attribute
	get json() {
		return JSON.parse(this.innerHTML);
	}
	get text() {
		return this.innerHTML;
	}
}
customElements.define("c-data", Data);
