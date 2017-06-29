import $ from "./minilib.js";

window.addEventListener("load", function(){
	var keymap = {
		ctrl : {
			"E" : goEditPage,
			//"M" : goMovePage,
			//"L" : goViewPage
		},
		alt : {
			//"P" : goPresentation
		}
	}

	document.body.on("keydown", function(e) {
		var keyCode = String.fromCharCode(e.keyCode);
		if(navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey){
			if(keymap.ctrl[keyCode]) {
				keymap.ctrl[keyCode]();
				e.preventDefault();
			}
		} else if(e.altKey) {
			if(keymap.alt[keyCode]){
				keymap.alt[keyCode]();
				e.preventDefault();
			}
		}
	}, false)

	function goEditPage() {
		location.href = "/!/edit" + location.pathname;
	}
})
