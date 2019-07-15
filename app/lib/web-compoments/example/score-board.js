import $ from "../minilib.module.js";

var template = $.template`
<style>
	:host {
	}
</style>
<slot name="header"></slot>
<ol>
	<template name="item">
		<li id="{{name}}">{{name}} - {{score}}</li>
	</template>
</ol>
`;

class ScoreBoard extends $.CustomElement {
	constructor() {
		super(template.content);

		this.on("attribute-changed", async (evt) => {

			await navigator.serviceWorker.register("./service-worker.js")
			await new Promise((resolve, reject) => { setTimeout(resolve, 1000)})
			console.log("load")
			await this.load(this.attr("from"))
		});

		$.event.on("score:updated", () => this.load(this.attr("from")));
	}
	static get observedAttributes() {
		return ["from"];
	}
	async load(url){
		if (!url) {
			console.warn("not enough information")
			return
		}
		var res = await $.request("GET", url)
		console.log("load complete");

		$.get(this["--shadow"], "ol").clear($.filters.exceptTemplate);

		var t = $.get(this["--shadow"], "template[name=item]")
		res.json.map((e) => {
			return $.render(t, e);
			// or
			// return $.create("li", {$text: `${e.name} - ${e.score}`, id: e.name});
		}).forEach((e) => {
			$.get(this["--shadow"], "ol").appendChild(e);
		});
	}
}
customElements.define("score-board", ScoreBoard);

