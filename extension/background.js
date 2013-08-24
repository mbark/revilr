function pageOnClick(info, tab) {
	query = "";
	query += "type=" + "page";
	query += "&url=" + info.pageUrl;
	query += "&title=" + tab.title;

	openPopup(query);
}

function imageOnClick(info, tab) {
	query = "";
	query += "type=" + "image";
	query += "&url=" + info.srcUrl;

	openPopup(query);
}

function selectionOnClick(info, tab) {
	query = "";
	query += "type=" + "selection";
	query += "&url=" + info.pageUrl;
	query += "&title=" + info.selectionText;

	openPopup(query);
}

function openPopup(query) {
	var url = "http://localhost:8080/revil";
	url += "?" + query;
	window.open(url);
}

function createContextMenu() {
	chrome.contextMenus.create({"title": "Revil this page!", "contexts":["page"], "onclick":pageOnClick});
	chrome.contextMenus.create({"title": "Revil image!", "contexts":["image"], "onclick":imageOnClick});
	chrome.contextMenus.create({"title": "Revil selection!", "contexts":["selection"], "onclick":selectionOnClick});
}

createContextMenu();