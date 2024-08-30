import * as $ from "bm.js/bm.module.js";

$.all("textarea[editor]").forEach(elem => elem.on("keydown", evt => {
    switch(evt.code) {
        case "Tab":
            evt.preventDefault();
            if (evt.shiftKey) {
                // TODO remove tab
            } else {
                let $textarea = evt.target;
                var start = $textarea.selectionStart;
                var end = $textarea.selectionEnd;
                var data = $textarea.value;

                if (end-start > 0 ) {
                    // mean selection is not empty
                    
                    let lines = [data.substring(0, start), data.substring(start, end), data.substring(end)];
                    
                    // 
                    let n = lines[0].lastIndexOf('\n');
                    console.log(n, start)
                    if (n > 0) {
                        lines[0] = [lines[0].substring(0, n+1), lines[0].substring(n+1)].join("\t");
                    }
                    
                    lines[1] = lines[1].split("\n").join("\n\t");

                    $textarea.value = lines.join("");
                    $textarea.selectionStart = $textarea.selectionEnd = start + 1;

                    console.log(lines);
                } else {
                    $textarea.value = data.substring(0, start) + "\t" + data.substring(end);
                    $textarea.selectionStart = $textarea.selectionEnd = start + 1;
                }
                // TODO indent selected line
            }
            return
        default:
            //console.log(evt);
    }
}))