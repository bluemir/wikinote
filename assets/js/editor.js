import * as $ from "bm.js/bm.module.js";

$.all("textarea[editor]").forEach(elem => elem.on("keydown", evt => {
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