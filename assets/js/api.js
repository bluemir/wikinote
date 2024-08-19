import * as $ from "bm.js/bm.module.js";

let cache = new Map();

export async function me() {
    if (cache.has("me")) {
        return cache.get("me");
    }

    let res = await $.request("GET", `/-/api/v1/me`, {withCredentials: true});

	cache.set("me", res.json);

    return res.json;
}


export async function can(verb, kind) {
    let key = `can/${verb}/${kind}`;

    if (cache.has(key)) {
        return cache.get(key)
    }

	let res = await $.request("GET", `/-/api/v1/can/${verb}/${kind}`);

    // not error
	cache.set(key, true);

    return true;
}

