import * as $ from "../../lib/bm.js/bm.module.js";
import {html, render} from 'lit-html';
import {Shortcut} from "../shortcut.js";

var sc = new Shortcut($.get("body"));
sc.add("ctrl+e", goToEdit);

function goToEdit(){
	location.assign(location.pathname+"?edit");
}
