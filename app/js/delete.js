import $ from "../lib/minilib.module.js";

var template = $.template`
<style>
.msg {
	display: none;
}
.msg.show {
	display: block;
}


</style>
If you want delete this page, put name of this page.
<form>
	<div><slot></slot></div>
	<input/>
	<button type="submit">Delete</button>
</form>
<p class="msg complete">
	Delete Complete. go to <a href="/">Front Page</a>
</p>
<p class="msg error">
	Something is wrong...
</p>
`

class DeleteForm  extends $.CustomElement {
	constructor() {
		super(template.content);

		$.get(this.shadow, "form").on("submit", (evt) => this.handleDelete(evt));
	}
	async handleDelete(evt) {
		evt.preventDefault();

		var url = location.pathname;
		try {
			var res = await $.request("DELETE", url, {
				header: {
					"X-Confirm": $.get(this.shadow, "input").value,
				},
			});
			if (res.statusCode < 300) {
				// show complete message
				$.get(this.shadow, ".msg.complete").classList.add("show")
			}
		} catch(e){
			// show error message
			$.get(this.shadow, ".msg.error").classList.add("show")
		}
	}
}
customElements.define("wikinote-delete-form", DeleteForm);
