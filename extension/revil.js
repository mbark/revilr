function renderText(url, comment) {
  	document.getElementById('urlField').value = url;
  	document.getElementById('commentField').value = comment;
}

function revilNow() {
	var type = getAndDeleteFromStorage('type');
	var targetUrl = "http://127.0.0.1:8080/revilr/" + type;

	var posting = $.post(targetUrl, {
		url: document.getElementById('urlField').value,
		c: document.getElementById('commentField').value
	});
	posting.done(function() {
		alert("Reviled!");
	});
	posting.fail(function() {
		alert("Failed to revil!");
	});
	posting.always(function() {
		window.close();
	});
}

function getAndDeleteFromStorage(item) {
	var stored = window.localStorage.getItem(item);
  	window.localStorage.removeItem(item);
  	return stored;
}

document.addEventListener("DOMContentLoaded", function () {
  	document.querySelector('button').addEventListener('click', revilNow);

  	var url = getAndDeleteFromStorage('url');
  	var comment = getAndDeleteFromStorage('comment');
  	renderText(url, comment);
});