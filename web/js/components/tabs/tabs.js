import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/static/css/common.css");

		:host {
		}
		header {
			display: flex;
		}
		::slotted(c-tab-panel) {
			display: none;
		}
		::slotted(c-tab-panel.selected) {
			display: block;
		}
	</style>
	<header @click=${evt => app.handleTabHeaderClick(evt)}>
		<slot name="header"></slot>
	</header>
	<slot name="panel"></slot>
`;

class Tabs extends $.CustomElement {
	constructor() {
		super();
	}
	static get observedAttributes() {
		return [ "selected" ];
	}
	async onConnected(){
		let role = this.attr("selected") || $.get(this, `c-tab-header`).attr("role");

		this.selected = role;
	}
	async onAttributeChanged(name, ov, nv) {
		switch(name) {
			case "selected":
				if (ov == nv) {
					return;
				}
				this.selected = nv;
				return
		}
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	async handleTabHeaderClick(evt) {
		let role = evt.target.attr("role");
		this.selected = role;
	}

	get selected() {
		return $.get(this, `c-tab-panel.show`).attr("role");
	}
	set selected(role) {
		this.attr("selected", role);
		$.all(this, `c-tab-header`).forEach(elem => elem.classList.remove("selected"));
		$.get(this, `c-tab-header[role=${role}]`).classList.add("selected");
		$.all(this, `c-tab-panel`).forEach(elem => elem.classList.remove("selected"));
		$.get(this, `c-tab-panel[role=${role}]`).classList.add("selected");
	}
}
customElements.define("c-tabs", Tabs);