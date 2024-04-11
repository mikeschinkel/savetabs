import {getApiServerUrl, getHttpOptions} from './api.js';

console.log("API host:", getApiServerUrl())

function getApiBrowseUiHtmlUrl() {
   return `${getApiServerUrl()}/html/browse`
}

const errorIcon = '<svg xmlns="http://www.w3.org/2000/svg" class="inline stroke-current flex-shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>';
document.addEventListener('DOMContentLoaded', function () {
   document.addEventListener('htmx:afterRequest', function (event) {
      if (event.detail.successful) {
         return
      }
      const xhr = event.detail.xhr
      const statusPanel = document.getElementById('status-panel');
      const errorNode = document.createElement('div');
      errorNode.innerHTML = xhr.responseText;
      const firstChild = statusPanel.firstChild;
      if (firstChild !== null) {
         statusPanel.insertBefore(errorNode.firstChild, firstChild);
      } else {
         statusPanel.appendChild(errorNode.firstChild);
      }
   });
});

// window.htmx.onLoad(function (target) {
//    // document.addEventListener('htmx:afterRequest', function (event) {
// });
// Attach to window for global availability

let clicked = false;
window.isBranchCollapsed = (id) => {
   let el = document.getElementById(id)
   if (el === null) {
      return true
   }
   const els = el.getElementsByTagName('ul')
   if (els.length === 0) {
      return true
   }
   const display = els[0].style.display
   if (display === 'none') {
      return true
   }
   if (display === '') {
      return true
   }
   return false;
}

class ErrorIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `
            <style>
                .icon-error {
                    stroke: currentColor; /* Use the text color of the parent element for the stroke */
                    flex-shrink: 0; /* Prevent the icon from shrinking in flex layouts */
                    height: 24px; /* Set height to 24 pixels */
                    width: 24px; /* Set width to 24 pixels */
                    fill: white;
                }
            </style>
            <svg xmlns="http://www.w3.org/2000/svg" class="icon-error" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
        `;
   }
}

customElements.define('error-icon', ErrorIcon);

class CloseIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `
            <style>
                .icon-close {
                    height: 24px;
                    width: 24px; 
                    fill: none; 
                    stroke: currentColor; /* Use the text color of the parent element for the stroke */
                    stroke-linecap: round;
                    stroke-linejoin: round;
                    stroke-width: 2px;
                }
            </style>
            <svg xmlns="http://www.w3.org/2000/svg" class="icon-close" viewBox="0 0 24 24">
                <path d="M6 18L18 6M6 6l12 12"></path>
            </svg>
        `;
   }
}

customElements.define('close-icon', CloseIcon);
