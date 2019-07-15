import $ from "../minilib.module.js";

var template = $.template`
<style>
	:host {
	}
	:host .msg {
		color: red;
	}
	:host .msg.hide {
		display : none;
	}
</style>
<form>
	<h3>register your score</h3>
	<p class="msg hide">
		something is wrong!
	</p>
	<input name="name" placeholder="name" />
	<input name="score" type="number" placeholder="score" />
	<button type="submit">Submit</button>
</form>

`;

class ScoreRegistryForm  extends $.CustomElement {
	constructor() {
		super(template.content);

		$.get(this["--shadow"], "form").on("submit", (evt) => this.submit(evt))
	}
	async submit(evt) {
		try {
			evt.preventDefault();
			this.disable();

			var to = this.attr("send-to");
			var type = this.attr("type") || "json";

			var shadow = this["--shadow"];

			var body;
			switch (type) {
				case "json":
					body = {
						name:  $.get(shadow, "input[name=name]").value,
						score: $.get(shadow, "input[name=score]").value,
					};
					break;
				default:
					console.error("unsupported type")
					$.event.fireEvent("message:warn", "unsupported type");
					return;
			}

			var res = await $.request("POST", to, {
				body: body,
			});

			this.fireEvent("updated", res.json);
			$.event.fireEvent("score:updated", res.json);
			$.event.fireEvent("message:info", "Score updated");
			this.enable();
		} catch (e) {
			console.log("submit failed", e);
			$.get(this["--shadow"], ".msg").classList.remove("hide");
			$.event.fireEvent("message:error", "submit failed");
		}
	}
	enable() {
		$.all(this["--shadow"], "input").forEach((e) => {e.value=""; e.attr("disabled", null) });
	}
	disable() {
		$.all(this["--shadow"], "input").forEach((e) => e.attr("disabled", ""));
	}
}
customElements.define("score-register-form", ScoreRegistryForm);

