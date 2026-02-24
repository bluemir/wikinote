import * as $ from "bm.js/bm.module.js";

$.all("button[type=button][method]").map(elem => elem.on("click", async evt => {
	evt.preventDefault();

	let $button  = evt.target;
	let method   = $button.attr("method") ? $button.attr("method").toLowerCase() : "get";
	let endpoint = $button.attr("action") ? $button.attr("action") : location;

	let res = await fetch(endpoint, {method: method});

	if (res.redirected) {
		location.href = res.url;
		return
	}
}));
