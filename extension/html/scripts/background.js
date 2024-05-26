import { apiPostLinksWithGroups } from './api.js';
import {
   getAllLinksWithGroups,
   setChromeObject,
   addRuntimeMessageListener,
} from './chromeUtils.js';

console.log("SaveTabs background.js loaded");
let intervalHandle;

setChromeObject(chrome)

/**
 * Collect all links periodically to submit, in inidividual ones failed to capture them
 * TODO: Refactor most logic into ./api.js
 */
function collectLinksWithGroups() {
   getAllLinksWithGroups()
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
            apiPostLinksWithGroups(links)
         })
         .catch(error => {
            // Handle errors
            console.log('Error: collectLinks(): ', error);
         });
   // clearInterval(intervalHandle);
}

// See https://medium.com/@otiai10/how-to-use-es6-import-with-chrome-extension-bd5217b9c978
intervalHandle = setInterval(collectLinksWithGroups, 60 * 1000);
collectLinksWithGroups()

/**
 * Capture message from content.js which replies with the Chrome Tab it was sent
 * from because Chrome won't let content scripts access tabs, Doh!
 */
addRuntimeMessageListener((message, sender, sendResponse) => {
   if (message.action !== "savetabs:getActiveTab") {
      return
   }
   sendResponse(sender.tab);
});

// addTabsCreatedListener(tab => {
//    putLinkByUrl(tab)
// });

// addTabsUpdatedListener((tabId, changeInfo, tab) => {
//    if (changeInfo.status !== 'complete') {
//       return
//    }
//    putLinkByUrl(tab)
//    console.log(`Tab updated: ${tab.title} - URL: ${tab.url}`);
// });

