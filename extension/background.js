function linkOnClick(info, tab) {
  	window.localStorage.setItem( 'type', 'page');
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
	chrome.windows.create({ url: popup, height: 280, type:"popup", focused:true });
}

function createContextMenu() {
chrome.contextMenus.create({"title": "Revil link!", "contexts":["link"], "onclick":linkOnClick});

chrome.contextMenus.create({"title": "Revil this page!", "contexts":["page"], "onclick":pageOnClick});

chrome.contextMenus.create({"title": "Revil image!", "contexts":["image"], "onclick":imageOnClick});

chrome.contextMenus.create({"title": "Revil selection!", "contexts":["selection"], "onclick":selectionOnClick});
}

createContextMenu();