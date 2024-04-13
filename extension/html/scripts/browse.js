import {getApiServerUrl, getHttpOptions} from './api.js';

console.log("API host:", getApiServerUrl())

function getApiBrowseUiHtmlUrl() {
   return `${getApiServerUrl()}/html/browse`
}

document.addEventListener('DOMContentLoaded', function () {
   document.addEventListener('htmx:afterRequest', function (event) {
      if (event.detail.successful) {
         return
      }
      const xhr = event.detail.xhr
      const statusPanel = document.getElementById('status-panel');
      const errorNode = document.createElement('div');
      errorNode.innerHTML = xhr.responseText;
      const isEmpty = errorNode.innerHTML.trim() === "";
      const noKids = errorNode.children.length === 0;
      if (isEmpty) {
         const div = document.createElement('div')
         div.innerText = `<div class="alert alert-error">Daemon API at ${getApiServerUrl()} appears unavailable.</div>`
         errorNode.innerHTML = div.textContent;
      } else if (noKids) {
         const span = document.createElement('span')
         span.innerText = errorNode.innerText
         errorNode.innerHTML = span.innerText
      }
      const firstChild = statusPanel.firstChild;
      if (firstChild !== null) {
         statusPanel.insertBefore(errorNode.firstChild, firstChild);
      } else {
         statusPanel.appendChild(errorNode.firstChild);
      }
   });
   document.addEventListener('htmx:targetError', function (event) {
      console.log('htmx:targetError: ', event)
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


// preventSummaryExpand stops propagation if the event targets 'SUMMARY'
window.preventSummaryExpand = (event) => {
   if (event.target.tagName !== 'SUMMARY') {
      return
   }
   event.preventDefault()
}

class ErrorIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `
            <style>
                .error-icon {
                    stroke: currentColor; /* Use the text color of the parent element for the stroke */
                    flex-shrink: 0; /* Prevent the icon from shrinking in flex layouts */
                    height: 24px; /* Set height to 24 pixels */
                    width: 24px; /* Set width to 24 pixels */
                    fill: white;
                }
            </style>
            <svg xmlns="http://www.w3.org/2000/svg" class="error-icon" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
        `;
   }
}

class CloseIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `
            <style>
                .close-icon {
                    height: 24px;
                    width: 24px; 
                    fill: none; 
                    stroke: currentColor; /* Use the text color of the parent element for the stroke */
                    stroke-linecap: round;
                    stroke-linejoin: round;
                    stroke-width: 2px;
                }
            </style>
            <svg xmlns="http://www.w3.org/2000/svg" class="close-icon" viewBox="0 0 24 24">
                <path d="M6 18L18 6M6 6l12 12"></path>
            </svg>
        `;
   }
}

class ExpandIcon extends HTMLElement {
   connectedCallback() {
      this.innerHTML = `
      <svg width="1rem" height="1rem" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M9 18l6-6-6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>`;
   }
}

class CollapseIcon extends HTMLElement {
   connectedCallback() {
      this.innerHTML = `
      <svg width="1rem" height="1rem" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>`;
   }
}

class BlankIcon extends HTMLElement {
   connectedCallback() {
      this.innerHTML = `
      <svg width="1rem" height="1rem" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path />
      </svg>`;
   }
}

customElements.define('error-icon', ErrorIcon);
customElements.define('close-icon', CloseIcon);
customElements.define('expand-icon', ExpandIcon);
customElements.define('collapse-icon', CollapseIcon);
customElements.define('blank-icon', BlankIcon);

document.addEventListener('alpine:init', () => {
   Alpine.data('collapsible', (initialState='collapsed') => ({
      state: initialState,
      expanded: function () {
         return this.state === 'expanded'
      },
      collapsed: function () {
         return this.state === 'collapsed'
      },
      collapse: function (ev) {
         this.state = 'collapsed';
         ev.stopPropagation();
      },
      expand: function (ev) {
         this.state = 'expanded';
         ev.stopPropagation();
      },
   }))
})

// Load Alpine.js after initial configurations are done
function loadAlpine() {
   const script = document.createElement('script');
   script.src = 'scripts/alpinejs@3.13.8.min.js';
   script.integrity='sha384-MGt/yQlIAvCVZEB4PNx8b9JxEfqFXemRJcpH6AIHAxDt1bRfYFeOnv3HJMW0LVD3';
   script.crossorigin='anonymous';
   script.defer = true;
   script.onload = () => {
      console.log('Alpine.js has loaded, initializing now...');
   };
   document.head.appendChild(script);
}
loadAlpine();
