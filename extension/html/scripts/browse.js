import {getApiServerUrl, getHttpOptions} from './api.js';
console.log("API host:", getApiServerUrl())

function getApiBrowseUiHtmlUrl() {
   return `${getApiServerUrl()}/html/browse`
}

// document.addEventListener('DOMContentLoaded', function () {
//    // alert("DOM Loaded")
//    const attr = 'hx-get';
//    let el = document.getElementById("browse-ui");
//    el.setAttribute(attr,`${getApiServerUrl()}/${el.getAttribute(attr)}`);
// });

// window.htmx.onLoad(function(target) {
//    alert("HTMX Loaded")
// });