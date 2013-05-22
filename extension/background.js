function allOnClick(info, tab) {
  alert(JSON.stringify(info) + "\n\n\n" + JSON.stringify(tab));
}

function linkOnClick(info, tab) {
  alert(info.linkUrl);
}

function pageOnClick(info, tab) {
	var url = info.pageUrl;
	var comment = tab.title;

	var popup = 'revil.html#' + url;
	chrome.windows.create({ url: popup, width: 450, height: 220 });
}

function imageOnClick(info, tab) {
  alert(info.srcUrl);
}

function selectionOnClick(info, tab) {
  alert(info.pageUrl + "\n\n\"" + info.selectionText + "\"");
}

function createContextMenu() {
var context = "link";
chrome.contextMenus.create({"title": "Revil " + context + "!", "contexts":[context], "onclick":linkOnClick});

var context = "page";
chrome.contextMenus.create({"title": "Revil " + context + "!", "contexts":[context], "onclick":pageOnClick});

var context = "image";
chrome.contextMenus.create({"title": "Revil " + context + "!", "contexts":[context], "onclick":imageOnClick});

var context = "selection";
chrome.contextMenus.create({"title": "Revil " + context + "!", "contexts":[context], "onclick":selectionOnClick});
}

createContextMenu();