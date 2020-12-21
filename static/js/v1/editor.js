import * as $ from "/!/static/lib/minilib.module.js";
import {html, render} from '/!/lib/lit-html/lit-html.js';
import {Shortcut} from "../shortcut.js";

var tmpl = (app) => html`
<style>
	@import url("/!/static/css/color.css");

	:host {
		display: block;
	}
	textarea {
		width: 100%;
		box-sizing: border-box;
		min-height: 40rem;
		tab-size: 4;
	}

	/* preview */
	.tabs .editor, .tabs .preview, .tabs .attribute {
		display: none;
	}
	.tabs[state=editor] .editor {
		display: block;
	}
	.tabs[state=preview] .preview {
		display: block;
	}

	.tabs[state=attribute] .attribute {
		display: block;
	}
	.tabs menu {
		margin: 0rem;
		margin-top: 1rem;
		padding: 0rem;

		border-bottom: 1px solid var(--button-background-color);
	}

	.tabs .btn[tab] {
		color: var(--button-font-color);
	}

	.tabs[state=editor]  .btn[tab=editor] {
		color: var(--button-selected-font-color);
	}
	.tabs[state=preview] .btn[tab=preview] {
		color: var(--button-selected-font-color);
	}
	.tabs[state=attribute] .btn[tab=attribute] {
		color: var(--button-selected-font-color);
	}


</style>
<section class="tabs marginer" state="editor">
	<menu class="tabs-header">
		<a class="btn" href="#" tab="editor"   >Edit</a>
		<a class="btn" href="#" tab="preview"  >Preview</a>
		<a class="btn" href="#" tab="attribute">Attribute</a>
	</menu>
	<section class="panel editor">
		<form method="POST" action="${location.pathname}">
			<textarea name="data">${app.data}</textarea>
			<input type="submit" value="Save"/>
		</form>
	</section>
	<section class="panel preview wiki-contents">
	</section>
	<section class="panel attribute">
		<script type="module" src="/!/static/js/kv-editor.js"></script>
		<kv-editor></kv-editor>
		<button x-func="attr-save">Save Attribute</button>
	</section>
</section>

`;

class WikinoteEditor extends $.CustomElement {
	constructor() {
		super();

		this.on("connected", () => this.onConnected())
	}

	onConnected() {
		console.log("connected")
		this.render();

		var sc = new Shortcut($.get("body"));
		sc.add("ctrl +space", e => this.previewToggle());
		sc.add("alt + space", e => this.previewToggle());
		sc.add("alt + .",     e => this.previewToggle());
		sc.add("alt + a",     e => this.attribute());

		var editorShotcut = new Shortcut($.get(this.shadow, "form"));
		editorShotcut.add("tab", e => this.addTab());
		editorShotcut.add("ctrl + s", e => this.save(e));
	}

	async render() {
		render(tmpl(this), this.shadow);
	}
	async save(evt) {
		let str = $.get(this.shadow, "textarea").value;
		let path = $.get(this.shadow, "form").attr("action");
		let res = await $.request("PUT", path, {
			body: str
		});

		// TODO show message
		console.log("saved");
	}
	addTab() {
		var $textarea = $.get(this.shadow, "textarea[name=data]");
		var start = $textarea.selectionStart;
		var end = $textarea.selectionEnd;
		var data = $textarea.value;

		$textarea.value = data.substring(0, start) + "\t" + data.substring(end);
		$textarea.selectionStart = $textarea.selectionEnd = start + 1;

	}
	// attribute
	get data() {
		return this.innerHTML
	}
	// event listener

}
customElements.define("wikinote-editor", WikinoteEditor);
