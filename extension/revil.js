function renderText(url, comment) {
  	document.getElementById('urlField').value = url;
  	document.getElementById('commentField').value = comment;
}

function revilNow() {
	var type = getFromStorage('type');
	var targetUrl = "http://127.0.0.1:8080/revilr/" + type;

	var posting = $.post(targetUrl, {
		url: document.getElementById('urlField').value,
		c: document.getElementById('commentField').value
	});
	posting.done(function() {
		alert("Revild " + getFromStorage('type') + "!");
	});
	posting.fail(function() {
		alert("Failed to revil!");
	});
	posting.always(function() {
		clearStorage();
		window.close();
	});
}

function getFromStorage(item) {
	var stored = window.localStorage.getItem(item);
  	return stored;
}

function clearStorage() {
  	window.localStorage.removeItem('type');
  	window.localStorage.removeItem('url');
  	window.localStorage.removeItem('comment');
}

document.addEventListener("DOMContentLoaded", function () {
  	document.querySelector('button').addEventListener('click', revilNow);

  	var url = getFromStorage('url');
  	var comment = getFromStorage('comment');
  	renderText(url, comment);
});