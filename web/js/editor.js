import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		textarea {
			width: 100%;
			min-height: 30rem;
			max-width: 100%;
			min-width: 100%;
		}
	</style>
	<c-tabs selected="editor">
		<c-tab-header slot="header" role="editor">Editor</c-tab-header>
		<c-tab-panel  slot="panel"  role="editor">
			<form method="POST" action="${location.pathname}">
				<textarea name="data"></textarea>
				<button>Save</button>
			</form>
		</c-tab-panel>
		<c-tab-header slot="header" role="preview">Preview</c-tab-header>
		<c-tab-panel  slot="panel"  role="preview">
			<!-- use slot? or import css -->
		</c-tab-panel>
	</c-tabs>
`;

class WikinoteEditor extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
		$.get(this.shadowRoot, "textarea").innerHTML = this.innerHTML;
	}
	// attribute
}
customElements.define("wikinote-editor", WikinoteEditor);
