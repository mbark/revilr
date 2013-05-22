/**
 * Renders the URL for the image, trimming if the length is too long.
 */
function renderText(url) {
  document.getElementById('urlField').value = url;
}

function revilNow() {
  alert("Sent:\n\n" + document.getElementById('urlField').value + "\n" + document.getElementById('commentField').value);
}

/**
 * Load the image in question and display it, along with its metadata.
 */
document.addEventListener("DOMContentLoaded", function () {
  // The URL of the image to load is passed on the URL fragment.
  var url = window.location.hash.substring(1);

  renderText(url);
});
