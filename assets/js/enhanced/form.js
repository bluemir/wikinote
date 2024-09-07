import * as $ from "bm.js/bm.module.js";

$.all(`form[method="delete" i], form[method="put" i]`).map(elem => elem.on("submit", async evt => {
	// hijack submit
	evt.preventDefault();

	let $form = evt.target;
	let method = $form.attr("method") ? $form.attr("method").toLowerCase() : "get";

	let data = new FormData($form);

	let res = await fetch($form.action, {method: method, body: data, redirect: "follow"});

	console.log("method=delete or method=put not support in valila, it will redirect 'GET' request")

	if (res.status >= 500) {
		console.error("error on request")
		return
	}

	// handle redirect
	if (res.redirected) {
		location.href = res.url;
		return
	}

	// if not,
	location.href = $form.action;
}))
