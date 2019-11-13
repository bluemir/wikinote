import $ from "../lib/minilib.module.js";
import {Shortcut} from "./shortcut.js";

class EditorController {
	constructor() {
		var sc = new Shortcut($.get("body"))
		sc.add("ctrl +space", e => this.previewToggle());
		sc.add("alt + space", e => this.previewToggle());
		sc.add("alt + .", e => this.previewToggle());

		var editorShotcut = new Shortcut($.get(".editor form"));
		editorShotcut.add("tab", e => this.addTab());
		editorShotcut.add("ctrl + s", e => this.save());
	}
	async previewOn() {
		var str = $.get("form textarea").value;
		var res = await $.request("POST", "/!/api/preview", {
			body: str
		})

		var $preview = $.get(".panel.preview");

		if ( res.statusCode>=200 && res.statusCode< 300) {
			$preview.innerHTML = res.text;
		} else {
			$preview.innerHTML = "Oops! error on get preview";
		}
		this.state = "preview";
	}
	previewOff() {
		this.state = "editor";
		$.get(".editor textarea").focus();
	}
	previewToggle() {
		if (this.state == "preview") {
			this.previewOff();
		} else {
			this.previewOn();
		}
	}
	get state() {
		return $.get(".tabs").attr("state")
	}
	set state(v) {
		$.get(".tabs").attr("state", v)
	}
	addTab() {
		var $textarea = $.get("textarea[name=data]");
		var start = $textarea.selectionStart;
		var end = $textarea.selectionEnd;
		var data = $textarea.value;

		$textarea.value = data.substring(0, start) + "\t" + data.substring(end);
		$textarea.selectionStart = $textarea.selectionEnd = start + 1;
	}
	async save() {
		var str = $.get(".tabs form textarea").value;
		var path = $.get(".tabs form").attr("action");
		$.request("PUT", path, {
			body: str
		});
	}
}

new EditorController();
