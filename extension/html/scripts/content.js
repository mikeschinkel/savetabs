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
   tab.html = document.documentElement.outerHTML
   console.log("Tab:", tab)
   apiPutLinkByTab(tab)
});

