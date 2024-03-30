import {getAllResources,setChromeObject,addRecentlySubmittedResources} from './chromeUtils.js';

const apiUrl = "http://localhost:8642/resources/with-groups"
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

function postResources(resources){
   if (resources.length===0) {
      return
   }
   const options = httpOptions('POST',resources)
   fetch(apiUrl, options)
      .then(async response => {
         console.dir(response)
         if (!response.ok) {
            // Handle HTTP errors (status codes other than 200-299)
            return response.text().then(errorResponse => {
               // Handle any errors (including HTTP errors)
               let status = `${response.status} - ${response.statusText}`
               console.log('Error: postResources(): ', status, errorResponse);
               throw new Error(`HTTP Error: ${status}. Details: ${JSON.stringify(errorResponse)}`);
            });
         }
         let text = await response.text()
         return text !== "" ? JSON.stringify(text) : {};
      })
      .then(data => {
         addRecentlySubmittedResources(resources)
         console.log('Response from server:', data); // Handle the successful response
      })
      .catch(error => {
         console.log('Error: postResources(): ', error); // Handle any errors (including HTTP errors)
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
         postResources(resources)
      })
      .catch(error => {
         // Handle errors
         console.log('Error: collectResources(): ', error);
      });
   // clearInterval(intervalHandle);
}

intervalHandle = setInterval(collectResources, 5000);

