function renderText(url) {
  document.getElementById('urlField').value = url;
}

function revilNow() {
  alert("Sent:\n\n" + document.getElementById('urlField').value + "\n" + document.getElementById('commentField').value);
}


document.addEventListener("DOMContentLoaded", function () {
  var url = window.location.hash.substring(1);
  document.querySelector('button').addEventListener('click', revilNow);
  renderText(url);
});
