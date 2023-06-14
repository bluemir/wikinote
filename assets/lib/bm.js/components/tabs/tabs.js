import * as $ from "../../bm.module.js";
import {html, render} from 'lit-html';
import {css} from "../common.js";
/*

<c-tabs selected="a">
	<c-tab-header slot="header" role="a">A Header</c-tab-header>
	<c-tab-panel  slot="panel"  role="a">
		A Contents
	</c-tab-panel>
	<c-tab-header slot="header" role="b">B Header</c-tab-header>
	<c-tab-panel  slot="panel"  role="b" @active="${evt => onActive(evt)}">
		B Contents
	</c-tab-panel>
</c-tabs>

*/

var tmpl = (app) => html`
	<style>
		${css}

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
		header {
			margin-bottom: 1rem;
			border-bottom: 1px solid gray;
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

		this.changePanel(role);
	}
	async onAttributeChanged(name, ov, nv) {
		switch(name) {
			case "selected":
				if (ov == nv) {
					return;
				}
				this.changePanel(nv);
				return
		}
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	async handleTabHeaderClick(evt) {
		let role = evt.target.attr("role");
		if (!role) {
			return
		}
		this.selected = role;
	}

	get selected() {
		return this.attr("selected");
	}
	set selected(role) {
		this.attr("selected", role);
	}
	changePanel(role) {
		$.all(this, `c-tab-header`).forEach(elem => elem.classList.remove("selected"));
		$.get(this, `c-tab-header[role=${role}]`).classList.add("selected");
		$.all(this, `c-tab-panel`).forEach(elem => elem.classList.remove("selected"));
		$.get(this, `c-tab-panel[role=${role}]`).classList.add("selected");

		$.all(this, `.selected`).forEach(e => {
			e.fireEvent("active")
		});
	}
}
customElements.define("c-tabs", Tabs);
