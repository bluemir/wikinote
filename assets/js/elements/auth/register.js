import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<style>
		@import url("/-/static/css/system/grid.css");

		:host {
			display: block;
		}
		c-input, section {
			display: block;	
		}
	</style>

	<form @submit="${ evt => elem.onSubmit(evt) }" >
		<c-input label="username"          name="username" type="text"     placeholder="your nickname. eg) bluemir" ></c-input>
		<c-input label="password"          name="password" type="password" placeholder="min-length: 6" ></c-input>
		<c-input label="password confirm"  name="confirm"  type="password" placeholder="same as password" ></c-input>
		<section>
			<input type="checkbox" id="terms"/>
			<label for="terms"> I read and agree to terms &amp; conditions.</label>
		</section>
		<c-button><button>Create Account</button></c-button>
	</form>
`;

class CustomElement extends $.CustomElement {
	constructor() {
		super();

	}
	async render() {
		render(tmpl(this), this.shadow);
	}

	async onSubmit(evt) {
		evt.preventDefault();

		let fd = new FormData($.get(this.shadowRoot, "form"));

		let res = await $.request("POST", `/-/api/v1/users`, {body:fd});

		location.href = "/-/welcome"
	}
}
customElements.define("wikinote-register", CustomElement);
