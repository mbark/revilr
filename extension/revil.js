function renderText(url, comment) {
  	document.getElementById('urlField').value = url;
  	document.getElementById('commentField').value = comment;
}

function revilNow() {
	var type = getAndDeleteFromStorage('type');
	var url = document.getElementById('urlField').value;
	var comment = document.getElementById('commentField').value;

	var targetUrl = "http://127.0.0.1:8080/revilr/" + type;
	var params = "url=" + url + "&c=" + comment;

	postRevilToServer(targetUrl, params);
}

function postRevilToServer(targetUrl, params) {
	var http = new XMLHttpRequest();
	http.open("POST", targetUrl, true);
	http.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	http.setRequestHeader("Content-length", params.length);
	http.setRequestHeader("Connection", "close");
	http.onreadystatechange = function() {//Call a function when the state changes.
		if(http.readyState == 4 && http.status == 200) {
			alert("Sent succesfully!");
			window.close();
		}
	}
	http.send(params);
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