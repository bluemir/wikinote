import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {css} from "common.js";


class CustomElement extends $.CustomElement {
	template() {
        return html`
            <style>
                ${css}

                :host {
                    display: block;
                }
            </style>
            <table>
				<thead>
					<th>Name</th>
					<th>Group</th>
				</thead>
				<tbody>
					${this.items?.map(item => html`
						<tr>
							<td>${item.name}</td>
							<td>${item.groups.join(",")}</td>
						</tr>
					`)}
				</tbody>
			</table>
        `;
    }
	constructor() {
		super();

		this.items = [];
	}
	async onConnected() {
		let res = await $.request("GET", this.attr("endpoint"));

		this.items = res.json.items;

		this.render();
	}
	async render() {
		render(this.template(), this.shadow);
	}
	// attribute
	get files() {
		return $.get(this, "c-data").json;
	}
}
customElements.define("roles-list", CustomElement);
