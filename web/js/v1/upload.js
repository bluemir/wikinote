import * as $ from "../../lib/bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");
	</style>
	<form @submit=${evt => app.onSubmit(evt)}>
		<dl>
			<dt>
				<label for="path">Path</label>
			</dt>
			<dd>
				<input name="path" value="${app.path}" />
			</dd>
			<dt>
				<label for="file"}>File</label>
			</dt>
			<dd>
				<input name="file" type="file" />
			</dd>
			<dd>
				<button type="submit">Save</button>
			</dd>
		</dl>
	</form>
`;

class WikinoteUpload extends $.CustomElement {
	constructor() {
		super();

	}

	onConnected() {
		this.render();

	}

	async render() {
		render(tmpl(this), this.shadow);
	}

	// attribute
	get path() {
		return this.attr("path");
	}

	// event listener
	async onSubmit(evt) {
		evt.preventDefault();

		let file = $.get(this.shadow, "input[type=file]").files[0];
		let path = $.get(this.shadow, "input[name=path]").value;

		let res = await $.request("PUT", path, {body: file});
	}
}
customElements.define("wikinote-upload", WikinoteUpload);
