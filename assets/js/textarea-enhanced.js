import * as $ from "bm.js/bm.module.js";

$.all("textarea[indent]").map(elem => elem.on("keydown", evt => {
    // handle indent, un-indent
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
            return
        default:
        //console.log(evt);
    }
}));

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

$.all("textarea[auto-resize]").map(elem => {
    // use `field-sizing: content` when available
    if (CSS.supports("field-sizing", "content")) {
        elem.style.fieldSizing = "content"
        return
    }

    // try old fashioned way.
    elem.style.height = `${elem.scrollHeight+2}px`;
    elem.on("input", evt => {
        // resize textarea
        let $textarea = evt.target;
        $textarea.style.height = `auto`; // it's magic, shrink area to fit contents
        $textarea.style.height = `${$textarea.scrollHeight+2}px`;
    })
});

$.all("textarea[submit-shortcut]").map(textarea => {
    textarea.on("keydown", async evt => {
        if (!(evt.code == "KeyS" && evt.ctrlKey)) {
            return // just skip
        }
        evt.preventDefault();

        let $form = evt.target.closest("form");
        let data = new FormData($form);
        
        let res = await $.request($form.method||"get", $form.action||location.pathname, {body: data});
        // TODO show message
    })
});