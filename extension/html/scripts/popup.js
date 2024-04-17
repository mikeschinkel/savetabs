import {getApiServerUrl,checkApiHealth} from './api.js';

// Function to load content from API
async function loadExtensionUiFromApi() {
   return new Promise((resolve, reject) => {
      let content;
      fetch(`${getApiServerUrl()}/ui/home`)
            .then(response => {
               content = response.text();
            })
            .catch(reason => {
               content = reason.text();
            })
      document.getElementById("content-section").innerHTML = content
   })
}

async function handleApiHealthCheck() {
   await checkApiHealth(isApiHealthy => {
      document.getElementById(getFocusSection(isApiHealthy)).style.display = "block";
      document.getElementById(getFocusSection(!isApiHealthy)).style.display = "none";
   });
}

function getFocusSection(isApiHealthy) {
   return isApiHealthy
      ? "extension-popup"
      : "no-daemon-warning";
}
let apiHealthCheckHandle;
// Call the function to handle API health check
document.addEventListener('DOMContentLoaded', function () {
   const _ = handleApiHealthCheck();
   apiHealthCheckHandle = setInterval(handleApiHealthCheck, 5000)
});


// Handle click event on popup.html buttons
// Replace with @click from AlpineJS
document.addEventListener('DOMContentLoaded', function () {
   document.addEventListener('click', function(event) {
      const target = event.target;
      if (target.tagName !== 'BUTTON') {
         return
      }
      if (!target.hasAttribute('data-href')) {
         return
      }
      chrome.tabs.create({url: target.getAttribute('data-href')});
   });
});


