import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {shortcut} from './shortcut.js';
import {css, events} from "common.js";

var tmpl = (elem) => html`
	<style>
		${css}

		textarea {
			padding: 1rem;
			min-height: 30rem;
			width: 100%;
			resize: vertical;
			box-sizing: border-box;
			tab-size: 4;
		}
	</style>
	<c-tabs selected="editor">
		<c-tab-header slot="header" role="editor">Editor</c-tab-header>
		<c-tab-panel  slot="panel"  role="editor">
			<form method="post" action="${location.pathname}">
				<textarea name="data" @keydown="${evt => elem.handleTextareaInput(evt)}">${elem.data}</textarea>
				<button>Save</button>
			</form>
		</c-tab-panel>
		<c-tab-header slot="header" role="preview">Preview</c-tab-header>
		<c-tab-panel  slot="panel"  role="preview" @active="${evt => elem.loadPreview(evt)}">
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
	onConnected() {
		$.get(this.shadowRoot, "textarea").focus();
		shortcut.add("ctrl+s", evt => this.save(evt))
	}
	async loadPreview(evt) {
		// TODO show 'now loading..."

		let data = $.get(this.shadowRoot, "form textarea")?.value;

		let res = await $.request("POST", "/-/api/preview", {
			data,
		});


		let elem = $.get(document, "[slot=preview]");

		elem.innerHTML = res.text;
	}
	handleTextareaInput(evt) {
		switch(evt.code) {
			case "Tab":
				evt.preventDefault();
				if (evt.shiftKey) {
					// TODO remove tab
				} else {
					this.addTab();
				}
				return
			default:
				//console.log(evt);
		}
	}
	addTab() {
		var $textarea = $.get(this.shadow, "textarea[name=data]");
		var start = $textarea.selectionStart;
		var end = $textarea.selectionEnd;
		var data = $textarea.value;

		$textarea.value = data.substring(0, start) + "\t" + data.substring(end);
		$textarea.selectionStart = $textarea.selectionEnd = start + 1;
	}
	async save(evt) {
		let str = $.get(this.shadowRoot, "textarea").value;
		let path = $.get(this.shadowRoot, "form").attr("action");
		let res = await $.request("PUT", path, {
			body: str
		});

		events.fireEvent("alert.info", {
			title: "saved",
		})
	}
	// attribute
	get data() {
		return $.get(this, "c-data").json?.data;
	}
}
customElements.define("wikinote-editor", WikinoteEditor);
