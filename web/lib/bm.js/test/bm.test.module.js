import * as $ from "../bm.module.js";
import "https://www.chaijs.com/chai.js"
import "https://unpkg.com/mocha/mocha.js"
// https://github.com/Tiliqua/assert-js

console.log(chai, mocha)
let assert = chai.assert;
let expect = chai.expect;

mocha.setup("tdd");
mocha.checkLeaks();

suite("AwaitEventTarget", () => {
	test("#test", async () => {
		let events = new $.AwaitEventTarget();

		events.on("page/next", async () => {
			await $.timeout(50);
			console.log("5 sec rule!")
		});
		events.on("page/next", async () => {
			await $.timeout(30);
			console.log("3 sec rule!")
		});


		let startTime = Date.now();

		await events.fireEvent("page/next");

		let endTime = Date.now();

		assert(endTime - startTime >= 50);
	})
});
suite("AwaitQueue", () => {
	test("#test", async () => {
		var q = new $.AwaitQueue();

		q.add(async () => {
			await $.timeout(50);
		});
		q.add(() => {
			return "return value"
		});

		let startTime = Date.now();
		// main event loop
		for(let func of q) {
			let ret = await func();
			if (ret == "return value"){
				return;
			}
		}
		let endTime = Date.now();

		assert(endTime - startTime >= 50);
	});
	test("#test promise", async () =>{
		var q = new $.AwaitQueue();

		var defer = $.defer();

		q.add(defer.promise);

		defer.resolve("end");

		let startTime = Date.now();
		// main event loop
		for(let v of q) {
			let ret = await v;
			if (ret == "end"){
				return;
			}
		}
		let endTime = Date.now();

		assert(endTime - startTime >= 50);
	})
	test("#anti pattern", async () => {
		var q = new $.AwaitQueue();

		q.add($.timeout(50));
		q.add($.timeout(30)); // not after 50ms, it just begin...

		let startTime = Date.now();
		// main event loop
		for(let v of q) {
			let ret = await v;
			if (ret == "end"){
				return;
			}
		}
		let endTime = Date.now();

		assert(endTime - startTime >= 50);
		assert(endTime - startTime < 80);
	})
});

mocha.run();
