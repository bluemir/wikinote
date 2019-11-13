import $ from "../lib/minilib.module.js";
import {Shortcut} from "./shortcut.js";

var sc = new Shortcut($.get("body"));
sc.add("ctrl + e", goToEdit)

function goToEdit(){
	location.href = location.pathname+"?edit";
}
