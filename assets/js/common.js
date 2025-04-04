import * as $ from "bm.js/bm.module.js";
import "bm.js/components/index.js";

let rev = $.get(`head script[type="importmap"]`).attr("rev");

export let css = `
@import url("/-/static/${rev}/css/element.css");
`;

export let events = new EventTarget();
