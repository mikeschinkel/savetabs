import {apiPutLinkByTab} from './api.js';
import {setChromeObject, lastError, sendRuntimeMessage} from './chromeUtils.js';

console.log("SaveTabs content.js loaded.");
setChromeObject(chrome)

console.log("Chrome:", chrome);

sendRuntimeMessage("", {action:"savetabs:getActiveTab"}, {}, (tab) => {
   const err = lastError()
   if (err) {
      console.log("ERROR:", err)
      return
   }
   console.log("Tab:", tab)
   apiPutLinkByTab(tab)
});

