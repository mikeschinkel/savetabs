import {getApiServerUrl,getHttpOptions} from './api.js';
import {getAllLinks,setChromeObject,addRecentlySubmittedLinks} from './chromeUtils.js';

console.log("SaveTabs loaded");
let intervalHandle;

setChromeObject(chrome)

function apiLinksWithGroupsEndpoint() {
   return `${getApiServerUrl()}/links/with-groups`
}

function postLinksWithGroups(lwgs){
   if (lwgs.length===0) {
      return
   }
   const options = getHttpOptions('POST',lwgs)
   fetch(apiLinksWithGroupsEndpoint(), options)
      .then(async response => {
         console.dir(response)
         if (!response.ok) {
            // Handle HTTP errors (status codes other than 200-299)
            return response.text().then(err => {
               // Handle any errors (including HTTP errors)
               let status = `${response.status} - ${response.statusText}`
               console.log('Error: postLinksWithGroups(): ', status, err);
               throw new Error(`HTTP Error: ${status}. Details: ${JSON.stringify(err)}`);
            });
         }
         let text = await response.text()
         return text !== "" ? JSON.stringify(text) : {};
      })
      .then(data => {
         addRecentlySubmittedLinks(lwgs)
         console.log('Response from server:', data); // Handle the successful response
      })
      .catch(err => {
         console.log('Error: postLinksWithGroups(): ', err); // Handle any errors (including HTTP errors)
      });
}

function collectLinks() {
   getAllLinks()
      .then(links => {
         if (links === undefined) {
            console.log('No new links to post');
            return
         }
         if (!Array.isArray(links)) {
            console.log('WARNING: Links not an array');
            return
         }
         if (links.length === 0) {
            console.log('No new links to post');
            return
         }
         console.log('Links:', links);
         postLinksWithGroups(links)
      })
      .catch(error => {
         // Handle errors
         console.log('Error: collectLinks(): ', error);
      });
   // clearInterval(intervalHandle);
}

intervalHandle = setInterval(collectLinks, 60 * 1000);
collectLinks()
