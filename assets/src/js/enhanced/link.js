import * as $ from "bm.js/bm.module.js";

$.all('[view=markdown] a[href^="http"]').map(elem => {
    if (elem.hostname !== location.hostname) {
        elem.attr("target", "_blank");
        elem.attr("rel", "noopener noreferrer");
    }
});
