import $ from "./minilib.js";

const KEYCODE_ALT = 18;
const KEYCODE_TAB = 8;

$.get("article form").on("submit", function(evt) {
	//console.log("hey!", this)
	//evt.preventDefault();
});

$.get("article form").on("keydown", function(evt) {
	var keyCode = String.fromCharCode(evt.keyCode);
	if(navigator.platform.match("Mac") ? evt.metaKey : evt.ctrlKey){
		if (keyCode == "S") {
			updateDocument();
			evt.preventDefault();
		}
	}
});
$.get("body").on("keydown", function(evt) {
	switch(evt.keyCode ) {
		case KEYCODE_ALT:
			previewOn();
			evt.preventDefault();
		return
	}
}, true)
$.get("body").on("keyup", function(evt) {
	switch(evt.keyCode ) {
		case KEYCODE_ALT:
			previewOff();
			evt.preventDefault();
		return
	}
}, true)
function updateDocument() {
	var str = $.get("article form textarea").value;
	var path = $.get("article form").getAttribute("action");
	$.request("PUT", path, {
		body: str
	});
}

var $preview = $.get(".panel.preview");
var $tabs  = $.get(".tab-control");
function previewOn(){
	var str = $.get("article form textarea").value;
	$.request("POST", "/!/api/preview", {
		body: str
	}).then(function (res) {
		if ( res.statusCode>=200 && res.statusCode< 300) {
			$preview.innerHTML = res.text;
		} else {
			$preview.innerHTML = "Oops! error on get preview";
		}
		// TOOD only class can handle
		$tabs.classList.add("preview");
		$tabs.classList.remove("editor");
	})
}
function previewOff(){
	$tabs.classList.remove("preview");
	$tabs.classList.add("editor");
}
