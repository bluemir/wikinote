import * as $ from "bm.js/bm.module.js";

function getIndentCharacter(attr) {
    switch(attr) {
        case "2space":
            return "  ";
        case "4space":
            return "    ";
        case "tab":
        default:
            return "\t";
    }
}

class EnhancedTextarea extends HTMLTextAreaElement{
	constructor() {
		super()
	}
	onConnected()  {
		if(this.hasAttribute("auto-resize")) {
			this.enableAutoResize();
		}

		if(this.hasAttribute("indent")) {
			this.enableIndent();
		}

		if(this.hasAttribute("submit-shortcut")) {
			this.enableSubmitShortcut();
		}
	}
	enableAutoResize() {
		if (CSS.supports("field-sizing", "content")) {
			this.style.fieldSizing = "content"
			return
		}

		// try old fashioned way.
		this.style.height = `${this.scrollHeight+2}px`;
		this.on("input", evt => {
			// resize textarea
			let $textarea = evt.target;
			$textarea.style.height = `auto`; // it's magic, shrink area to fit contents
			$textarea.style.height = `${$textarea.scrollHeight+2}px`;
		})
	}
	enableIndent() {
		this.on("keydown", this.#indent);

		if(this.hasAttribute("tab-size")) {
			this.style.tabSize = `${this.attr("tab-size")}ch`
		}
	}
	#indent(evt) {
		switch(evt.code) {
			case "Tab":
				evt.preventDefault();
				let $textarea = evt.target;
				let start = $textarea.selectionStart;
				let end = $textarea.selectionEnd;
				let data = $textarea.value;
				let indent = getIndentCharacter($textarea.attr("indent"))

				if (evt.shiftKey) {
					// un-tab

					let n = data.substring(0, start).lastIndexOf("\n")+1;

					let sections = [data.substring(0, n), data.substring(n, end), data.substring(end)];
					sections[1] = sections[1].split('\n').map(line => line.startsWith(indent)?line.substring(indent.length): line).join('\n');

					$textarea.value = sections.join("");

					$textarea.selectionStart = start > 0 ? start-1: 0;
					$textarea.selectionEnd   = sections[0].length + sections[1].length;
				} else {
					// tab
					// if (end-start > 0 ) { }// mean selection is not empty

					let n = data.substring(0, start).lastIndexOf("\n")+1;
					let sections = [data.substring(0, n), data.substring(n, end), data.substring(end)];

					sections[1] = sections[1].split('\n').map(line => indent + line).join('\n');

					$textarea.value = sections.join("");

					$textarea.selectionStart = start + indent.length;
					$textarea.selectionEnd   = sections[0].length + sections[1].length;
				}
				return;
			case "Enter":
				{
					let $textarea = evt.target;
					let start = $textarea.selectionStart;
					let end = $textarea.selectionEnd;
					let data = $textarea.value;

					if (end - start > 0) {
						return; // skip. it has selection
					}

					evt.preventDefault();

					// insert newline & indent
					let n = data.substring(0, start).lastIndexOf("\n")+1;

					let lastLine = data.substring(n, start);

					let indent = lastLine; // if empty line, use whole line.
					let matched = lastLine.match(/[^\s]/);
					if (matched){
						indent = data.substring(n, n + lastLine.match(/[^\s]/).index);
					}

					let arr = [data.substring(0, start), "\n", indent, data.substring(start)];

					$textarea.value = arr.join("");
					$textarea.selectionStart = $textarea.selectionEnd = arr[0].length + arr[1].length + arr[2].length;

					return;
				}
			default:
				//console.log(evt);
		}
	}
	enableSubmitShortcut() {
		this.on("keydown", this.#submitShortcut)
	}
	async #submitShortcut(evt) {
		if (!(evt.code == "KeyS" && evt.ctrlKey)) {
			return // just skip
		}
		evt.preventDefault();

		// evt.target.closest("form").submit();
		let $form = this.closest("form");

		let data = new FormData($form);
		let res = await $.request($form.attr("method")||$form.method, $form.action||location.pathname, {body: data});
		// TODO show message
	}
}
customElements.define("enhanced-textarea", EnhancedTextarea, {extends: "textarea"});

