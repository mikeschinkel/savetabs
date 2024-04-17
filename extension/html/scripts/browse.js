import {getApiServerUrl} from './api.js';

console.log("SaveTabs daemon:", getApiServerUrl())

function getApiBrowseUiHtmlUrl() {
   return `${getApiServerUrl()}/html/browse`
}

function getErrorNodeFromXHR(xhr) {
   const errorWrapper = document.createElement('div');
   let errorMsg = xhr.responseText;
   if (errorMsg.trim() === "") {
      errorWrapper.innerHTML = `<div class="alert alert-error">Daemon API at ${getApiServerUrl()} appears unavailable.</div>`;
      return errorWrapper
   }
   errorWrapper.innerHTML = errorMsg
   if (errorWrapper.children.length === 0) {
      const div = document.createElement('<div>')
      div.innerText = errorMsg
      errorWrapper.innerHTML = div.innerHTML
   }
   return errorWrapper.firstChild
}

function getStatusPanel() {
   let sp = document.getElementById('status-panel');
   if (sp===null){
      sp = document.createElement('div');
      document.body.appendChild(document.createElement('span'));
      document.body.insertBefore(sp, document.body.firstChild);
   }
   return sp;
}

document.addEventListener('DOMContentLoaded', function () {
   document.addEventListener('htmx:afterRequest', function (event) {
      if (event.detail.successful) {
         return
      }
      const errorNode = getErrorNodeFromXHR(event.detail.xhr);
      const statusPanel = getStatusPanel();
      const targetElem = statusPanel.firstChild;
      if (targetElem !== null) {
         statusPanel.insertBefore(errorNode, targetElem);
      } else {
         statusPanel.appendChild(errorNode);
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

document.addEventListener('alpine:init', () => {
   Alpine.data('preventable', () => ({
      preventExpandOnIconClick: function (event) {
         if (['svg','path'].includes(event.target.tagName.toLowerCase())) {
            return;
         }
         event.preventDefault()
      }
   }))
   Alpine.data('collapsible', () => ({
      state: 'collapsed',
      expanded: function () {
         return this.state === 'expanded'
      },
      collapsed: function () {
         return this.state === 'collapsed'
      },
      toggle: function () {
         switch (this.state) {
            case 'collapsed':
               this.state = 'expanded';
               break;
            case 'expanded':
               this.state = 'collapsed';
               break;
            default:
               alert(`Unexpected collapsible.state: ${this.state}`);
         }
      },
   }))
})

// Load Alpine.js after initial configurations are done
function loadAlpine() {
   const script = document.createElement('script');
   script.src = 'scripts/alpinejs-csp@3.13.8.min.js';
   script.integrity = 'sha384-ngWUWXGZwWaMILMRRZR1rmY/dFdl79tJTC1pJD+w4Ca0uyMJOf0p+L5C1fkWht0l';
   script.crossorigin = 'anonymous';
   script.defer = true;
   script.onload = () => {
      console.log('Alpine.js has loaded, initializing now...');
   };
   document.head.appendChild(script);
}
loadAlpine();
