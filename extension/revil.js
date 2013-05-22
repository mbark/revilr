function renderText(url) {
  document.getElementById('urlField').value = url;
}

function revilNow() {
	var http = new XMLHttpRequest();
	var targetUrl = "http://127.0.0.1:8080/revilr/page";
	var params = "url=" + document.getElementById('urlField').value + "&c=" + document.getElementById('commentField').value;
	http.open("POST", targetUrl, true);
	http.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	http.setRequestHeader("Content-length", params.length);
	http.setRequestHeader("Connection", "close");
	http.onreadystatechange = function() {//Call a function when the state changes.
		if(http.readyState == 4 && http.status == 200) {
			alert("Sent succesfully!");
		}
	}
	http.send(params);
}


document.addEventListener("DOMContentLoaded", function () {
  var url = window.location.hash.substring(1);
  document.querySelector('button').addEventListener('click', revilNow);
  renderText(url);
});
