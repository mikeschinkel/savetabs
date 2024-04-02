import {getApiServerUrl} from './api.js';
import {getAllResources,setChromeObject,addRecentlySubmittedResources} from './chromeUtils.js';

console.log("SaveTabs loaded");
let intervalHandle;

setChromeObject(chrome)

function httpOptions(method,data) {
   return {
      method: method,
      headers: {
         'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
   }
}
function apiResourcesWithGroupsEndpoint() {
   return `${getApiServerUrl()}/resources/with-groups`
}

function postResourcesWithGroups(rwgs){
   if (rwgs.length===0) {
      return
   }
   const options = httpOptions('POST',rwgs)
   fetch(apiResourcesWithGroupsEndpoint(), options)
      .then(async response => {
         console.dir(response)
         if (!response.ok) {
            // Handle HTTP errors (status codes other than 200-299)
            return response.text().then(err => {
               // Handle any errors (including HTTP errors)
               let status = `${response.status} - ${response.statusText}`
               console.log('Error: postResourcesWithGroups(): ', status, err);
               throw new Error(`HTTP Error: ${status}. Details: ${JSON.stringify(err)}`);
            });
         }
         let text = await response.text()
         return text !== "" ? JSON.stringify(text) : {};
      })
      .then(data => {
         addRecentlySubmittedResources(rwgs)
         console.log('Response from server:', data); // Handle the successful response
      })
      .catch(err => {
         console.log('Error: postResourcesWithGroups(): ', err); // Handle any errors (including HTTP errors)
      });
}

function collectResources() {
   getAllResources()
      .then(resources => {
         if (resources === undefined) {
            console.log('No new resources to post');
            return
         }
         if (!Array.isArray(resources)) {
            console.log('WARNING: resources not an array');
            return
         }
         if (resources.length === 0) {
            console.log('No new resources to post');
            return
         }
         console.log('Resources:', resources);
         postResourcesWithGroups(resources)
      })
      .catch(error => {
         // Handle errors
         console.log('Error: collectResources(): ', error);
      });
   // clearInterval(intervalHandle);
}

intervalHandle = setInterval(collectResources, 5000);

