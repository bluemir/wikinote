import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/!/static/css/color.css");

		table {
			border-collapse: collapse;
		}
		table th, table td {
			border: 1px solid gray;
		}
		table th {
			text-align: right;
			font-family: monospace;
			font-size: 1rem;
		}
	</style>
	<table>
		<thead>
		</thead>
		<tbody>
			${Object.entries(app.data|| {}).map(([filename, matches]) => html`
				<tr colspan="2">
					${filename}
				</tr>
				${matches.map((m) => html`
					<tr>
						<th>${m.line}</th>
						<td>${m.text}</td>
					</tr>
				`)}
			`)}
		</tbody>
	</table>
`;

class WikinoteSearch extends $.CustomElement {
	constructor() {
		super();
	}
	async render() {
		render(tmpl(this), this.shadow);
		let data = $.get(this, "c-data").json;
		this.data = data;
		console.log(data);
	}
	// attribute
}
customElements.define("wikinote-search", WikinoteSearch);
