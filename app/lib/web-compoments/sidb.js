// simple IndexedDB lib
//
function openDB({name, version, update} = {update:()=>{}}){
	var request = indexedDB.open(name);
	request.addEventListener("upgradeneeded", (evt) => {
		var db = evt.target.result;
		update(db)
	})
	return new Promise((resolve, reject) => {
		request.onsuccess = (evt) => {
			var db = evt.target.result;

			resolve(new SIDB(db))
		}
		request.onerror = reject;
	})
}
class SIDB {
	constructor(db) {
		this.inner = db
	}
	transaction(stores, mode) {
		var tx = this.inner.transaction(stores, mode)
		return new SIDBTransaction(tx)
	}
}
class SIDBTransaction{
	constructor(tx) {
		this.inner = tx

		this.complete = new Promise((resolve, reject) => {
			tx.oncomplete = resolve
			tx.onerror = reject
		})
	}
	store(name) {
		var store = this.inner.objectStore(name)
		return new SIDBObjectStore(store)
	}
}
class SIDBObjectStore {
	constructor(os) {
		this.inner = os
	}
	async get(key) {
		var req = this.inner.get(key)
		return new Promise((resolve, reject) => {
			req.onsuccess = (evt) => {
				resolve(evt.target.result)
			};
			req.onerror = reject;
		})
	}
	async put(data) {
		var req = this.inner.put(data)
		return new Promise((resolve, reject) => {
			req.onsuccess = (evt) => {
				resolve(evt.target.result)
			};
			req.onerror = reject;
		})
	}
	async list() {
		var ctx = await new Promise(resolve =>{
			var req = this.inner.openCursor();
			req.onsuccess = (evt) => { var cur = evt.target.result ; resolve(new SIDBCursor(req, cur)); }
		});
		return ctx
	}
}

class SIDBCursor {
	// TODO async?
	constructor(req, cursor) {
		this.req = req;
		this.inner = cursor;
	}
	*[Symbol.iterator](){
		while (this.inner) {
			yield this._wait_next(this.inner.value);
		}
	}
	async _wait_next(value) {
		this.inner = await new Promise(resolve => {
			this.req.onsuccess = evt => resolve(evt.target.result);
			this.inner.continue();
		})
		return value;
	}
}
/*
class EventIter {
	constructor(obj){
		this.target = obj
	}
	async *[Symbol.asyncIterator](){
		while(true) {
			yield await new Promise(resolve => {
				obj.onclick = function(evt) {resolve(evt)};
			});
		}
	}
}
*/
