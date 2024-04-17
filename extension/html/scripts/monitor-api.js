import {getApiServerUrl,checkApiHealth} from './api.js';

// Function to load content from API
// async function loadExtensionUiFromApi() {
//    return new Promise((resolve, reject) => {
//       let content;
//       fetch(`${getApiServerUrl()}/ui/home`)
//             .then(response => {
//                content = response.text();
//             })
//             .catch(reason => {
//                content = reason.text();
//             })
//       document.getElementById("content-section").innerHTML = content
//    })
// }

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
