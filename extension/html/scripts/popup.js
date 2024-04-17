import {checkApiHealth} from "./api.js";
import {} from "./alpine-loader.js";

// Handle click event on popup.html buttons
// Replace with @click from AlpineJS
document.addEventListener('DOMContentLoaded', function () {
   document.addEventListener('click', function (event) {
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


let apiHealthCheckHandle;

// Call the function to handle API health check
document.addEventListener('DOMContentLoaded', function () {
   // Call immediately to present popup
   const _ = handleApiHealthCheck()
   // Keep calling every 5 seconds to monitor the API
   apiHealthCheckHandle = setInterval(handleApiHealthCheck, 5000)
});

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

