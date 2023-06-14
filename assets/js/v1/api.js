import * as $ from "bm.js/bm.module.js";


export async function me() {
	let res = await $.request("GET", "/!/auth/me");
	return res.json;
}
