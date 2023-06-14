import * as $ from "bm.js/bm.module.js";

export async function can(verb, kind) {
	let res = await $.request("GET", `/-/api/auth/can/${verb}/${kind}`);

	return res.json;
}

