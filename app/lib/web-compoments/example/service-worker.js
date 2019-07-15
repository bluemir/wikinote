// work like server
//importScripts("../sidb.js");
importScripts("../sidb.js")

var db
this.addEventListener("install", (evt) => {
	evt.waitUntil(openDB({
			name: "rank-store",
			update: async (db) => await db.createObjectStore("rank", { keyPath: "name" }),
	}).then(async function(r) {
		db = r
		console.log(r)
		//await caches.open("data");
		var res = await fetch("./data.json");
		var data = await res.json();

		var tx = db.transaction(["rank"], "readwrite");
		var store = tx.store("rank");

		data.forEach(async (e) => {
			await store.put(e);
		})
		await tx.complete
	}));
})
this.addEventListener("activate", evt => {
	evt.waitUntil(clients.claim().then(function(){ console.log("activate");}))
})
this.addEventListener("fetch", function(evt) {
	if (evt.request.method == "POST") {
		console.log("POST...")
		evt.respondWith(saveScore(evt.request));
		return
	}
	if (evt.request.method == "GET" && evt.request.url.includes("data.json")) {
		console.log("GET data.json")
		evt.respondWith(loadScore(evt.request));
		return
	}
});


async function loadScore() {
	console.log(db);
	var tx = db.transaction(["rank"], "readwrite")
	var store = tx.store("rank");

	var result = await store.list();

	var rank = [];
	for await (var s of result) {
		rank.push(s)
	}

	await tx.complete

	rank.sort((a, b) => b.score - a.score)

	var res = new Response(
		JSON.stringify(rank),
		{ headers: { "Content-Type": "application/json" }},
	);
	return res
}

async function saveScore(req) {
	var s = await req.json();

	var tx = db.transaction(["rank"], "readwrite")
	var store = tx.store("rank");
	await store.put(s);

	var res = new Response(
		JSON.stringify(s),
		{ headers: { "Content-Type": "application/json" }},
	);
	return res;
}
