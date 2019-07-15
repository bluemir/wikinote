import $ from "../minilib.module.js";

var template = $.template`
<style>
	:host {
	}
	:host .info {
		color: blue;
	}
	:host .warn {
		color: yellow;
	}
	:host .error {
		color: red;
	}

</style>
<div>
	<template name="item">
		<p class="{{level}}">{{text}}</p>
	</template>
</div>
`;

class MessagePanel extends $.CustomElement {
	constructor() {
		super(template.content);

		$.event.on("message:info",  (msg) => this.add("info",  msg.detail));
		$.event.on("message:warn",  (msg) => this.add("warn",  msg.detail));
		$.event.on("message:error", (msg) => this.add("error", msg.detail));
	}
	add(level, text) {
		var t = $.get(this["--shadow"], "template[name=item]")
		$.get(this["--shadow"], "div").appendChild($.render(t, { level, text }));
	}
}
customElements.define("message-panel", MessagePanel);
