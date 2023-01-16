import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		textarea {
			min-height: 30rem;
			width: 100%;
			resize: vertical;
			box-sizing: border-box;
		}
	</style>
	<c-tabs selected="editor">
		<c-tab-header slot="header" role="editor">Editor</c-tab-header>
		<c-tab-panel  slot="panel"  role="editor">
			<form method="POST" action="${location.pathname}">
				<textarea name="data">${app.data}</textarea>
				<button>Save</button>
			</form>
		</c-tab-panel>
		<c-tab-header slot="header" role="preview">Preview</c-tab-header>
		<c-tab-panel  slot="panel"  role="preview" @active="${evt => app.loadPreview(evt)}">
			<!-- use slot? or import css -->
			<slot name="preview"></slot>
		</c-tab-panel>
	</c-tabs>
`;

class WikinoteEditor extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
	}
	async loadPreview(evt) {
		// TODO show 'now loading..."

		let data = $.get(this.shadowRoot, "form textarea")?.value;

		let res = await $.request("POST", "/!/api/preview", {
			data,
		});


		let elem = $.get(document, "[slot=preview]");

		elem.innerHTML = res.text;
	}
	// attribute
	get data() {
		console.log($.get(this, "c-data").text);
		return $.get(this, "c-data").text;
	}
}
customElements.define("wikinote-editor", WikinoteEditor);
