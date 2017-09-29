import $ from "./minilib.js";

$.get("#dropzone").on("drop", function(event) {
	Array.prototype.forEach.call(event.dataTransfer.files, uploadFile)
	event.preventDefault();
})
$.get("#dropzone").on("dragover", function(event) {
	$.get("#dropzone").style.color = "red";
	//ev.dataTransfer.dropEffect = "move"
	event.preventDefault();
})
async function uploadFile(file){
	console.log("filename: ", file.name)
	var basepath = $.get("#dropzone").getAttribute("x-path")

	var xhr = new XMLHttpRequest();

	xhr.upload.addEventListener("progress", function (e) {
		if (e.lengthComputable) {
			console.log(e.loaded / e.total);
		}
	});

	xhr.upload.addEventListener("load", function () {
		console.log("uploaded");
	});

	var filepath = basepath.slice(0, -3)+"/" +file.name;
	console.log(filepath);
	xhr.open("PUT", filepath);
	xhr.overrideMimeType(file.type);
	xhr.send(file);
}
