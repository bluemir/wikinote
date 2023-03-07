import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (app) => html`
	<style>
		@import url("/-/static/css/color.css");

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
		tbody a {
			color: var(--link-fg-color);
			text-decoration: none;
		}
		tbody a:hover {
			text-decoration: underline;
		}
	</style>
	<table>
		<thead>
		</thead>
		<tbody>
			${Object.entries(app.data|| {}).map(([filename, matches]) => html`
				<tr>
					<a href="${filename}">${filename}</a>
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
		this.data = $.get(this, "c-data").json;

		render(tmpl(this), this.shadow);
	}
	// attribute
}
customElements.define("wikinote-search", WikinoteSearch);
