import {getApiServerUrl} from './api.js';
const apiServerUrl = "http://localhost:8642"

async function checkApiHealth(callback) {
   const response = await fetch(`${apiServerUrl}/healthz`).catch( (_)=>{return {ok:false}})
   callback(response.ok)
}

// Function to load content from API
async function loadExtensionUiFromApi() {
   return new Promise((resolve, reject) => {
      let content;
      fetch(`${apiServerUrl}/ui/home`)
            .then(response => {
               content = response.text();
            })
            .catch(reason => {
               content = reason.text();
            })
      document.getElementById("content-section").innerHTML = content
   })
}
function getFocusSection(isApiHealthy) {
   return isApiHealthy
      ? "extension-popup"
      : "no-daemon-warning";
}
async function handleApiHealthCheck() {
   await checkApiHealth(isApiHealthy => {
      document.getElementById(getFocusSection(isApiHealthy)).style.display = "block";
      document.getElementById(getFocusSection(!isApiHealthy)).style.display = "none";
   });
}

// Call the function to handle API health check
document.addEventListener('DOMContentLoaded', function () {
   _ = handleApiHealthCheck();
   intervalHandle = setInterval(handleApiHealthCheck, 5000)

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


