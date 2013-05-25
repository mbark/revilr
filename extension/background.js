function linkOnClick(info, tab) {
  	window.localStorage.setItem( 'type', 'link');
	window.localStorage.setItem( 'url', info.linkUrl);

	openPopup();
}

function pageOnClick(info, tab) {
	window.localStorage.setItem( 'type', 'page');
	window.localStorage.setItem( 'url', info.pageUrl);
	window.localStorage.setItem( 'comment', tab.title);

	openPopup();
}

function imageOnClick(info, tab) {
  	window.localStorage.setItem( 'type', 'image');
	window.localStorage.setItem( 'url', info.srcUrl);

	openPopup();
}

function selectionOnClick(info, tab) {
	window.localStorage.setItem( 'type', 'selection');
	window.localStorage.setItem( 'url', info.pageUrl);
	window.localStorage.setItem( 'comment', info.selectionText);

	openPopup();
}

function openPopup() {
	var popup = 'revil.html';
	chrome.windows.create({ url: popup, width: 450, height: 155, type:"popup" });
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