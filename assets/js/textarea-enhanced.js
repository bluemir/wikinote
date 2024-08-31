import * as $ from "bm.js/bm.module.js";

$.all("textarea[indent-by-tab]").map(elem => elem.on("keydown", evt => {
    // handle indent, un-indent
    switch(evt.code) {
        case "Tab":
            evt.preventDefault();
            let $textarea = evt.target;
            let start = $textarea.selectionStart;
            let end = $textarea.selectionEnd;
            let data = $textarea.value;

            if (evt.shiftKey) {
                // un-tab

                let n = data.substring(0, start).lastIndexOf("\n")+1;

                let lines = [data.substring(0, n), data.substring(n, end), data.substring(end)];
                lines[1] = lines[1].split('\n').map(line => line.startsWith('\t')?line.substring(1): line).join('\n');

                $textarea.value = lines.join("");

                $textarea.selectionStart = start > 0 ? start-1: 0;
                $textarea.selectionEnd   = lines[0].length + lines[1].length;
            } else {
                // if (end-start > 0 ) { }// mean selection is not empty

                let n = data.substring(0, start).lastIndexOf("\n")+1;
                let lines = [data.substring(0, n), data.substring(n, end), data.substring(end)];

                lines[1] = lines[1].split('\n').map(line => "\t" + line).join('\n');

                $textarea.value = lines.join("");

                $textarea.selectionStart = start + 1;
                $textarea.selectionEnd   = lines[0].length + lines[1].length;
            }
            return
        default:
        //console.log(evt);
    }
}))

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
})