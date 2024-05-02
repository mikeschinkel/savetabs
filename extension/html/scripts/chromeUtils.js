// chromeUtils.js
const chromeUrlsRegex = /^(chrome(-\w+)?):\/\//;
let recentlySubmittedLinks = {};
let chrome = {
   windows: {
      getAll: (tabGroupId, callback) => throwNotSetError(),
   },
   tabGroups: {
      get: (tabGroupId, callback) => throwNotSetError(),
   },
   tabs: {
      get: (tabId, fn) => throwNotSetError(),
      query: (queryInfo, callback) => throwNotSetError(),
      sendMessage: (tabId,msg,opts,cb) => throwNotSetError(),
      onCreated: {
         addListener: (fn) => throwNotSetError(),
      },
      onActivated: {
         addListener: (fn) => throwNotSetError(),
      },
      onUpdated: {
         addListener: (fn) => throwNotSetError(),
      },
   },
   runtime: {
      lastError: {
         message: "chrome object not set with setChromeObject()"
      },
      onMessage: {
         addListener: (fn) => throwNotSetError(),
      },
   },
   scripting: {
      executeScript: (obj,fn)=> throwNotSetError(),
   },
};

export function setChromeObject(chromeObj) {
   chrome = chromeObj;
}

export function getTabGroupLinks(tabGroup) {
   return new Promise((resolve, reject) => {
      chrome.tabs.query({ groupId: tabGroup.id }, tabs => {
         if (chrome.runtime.lastError) {
            reject(chrome.runtime.lastError.message);
         } else {
            resolve(tabs
               .filter(tab =>  !rejectUrl(tab.url,tabGroup.title) )
               .map(tab => {
                  return {
                     original_url: tab.url,
                     title: tab.title,
                     groupId: tabGroup.id,
                     groupType: 'tab-group',
                     group: tabGroup.title,
                  };
               })
            );
         }
      });
   });
}

export function getAllTabGroupsLinks(tabGroups) {
    // Create a promise to retrieve tabs for each tab group
    const tabGroupPromises = tabGroups.map(tabGroup => {
        return getTabGroupLinks(tabGroup);
    });

    // Create a promise to retrieve tabs without a tab group
    const noGroupPromise = getTabGroupLinks({
       id: -1,
       title: '<none>',
    })

    // Combine promises and flatten results
    const allPromises = tabGroupPromises.concat(noGroupPromise);
    return Promise.all(allPromises)
        .then(results => {
           return results.flat();
        })
        .catch(error => {
            console.log('Error:', error);
            throw error;
        });
}

export function getAllTabGroups() {
   return new Promise((resolve, reject) => {
      chrome.tabGroups.query({}, (tabGroups) => {
         if (chrome.runtime.lastError) {
            reject(chrome.runtime.lastError.message);
         } else {
            resolve(tabGroups);
         }
      });
   });
}

export async function getAllLinksWithGroups() {
   try {
      const tabGroups = await getAllTabGroups();
      return await getAllTabGroupsLinks(tabGroups)
   } catch (error) {
      console.log('Error: getAllLinks(): ', error);
      throw error;
   }
}

export function addRecentlySubmittedLinks(links) {
   let d = new Date();
   links.forEach(r => {
      recentlySubmittedLinks[urlKey(r.original_url,r.group)] = d;
   })
}

export function addRuntimeMessageListener(fn) {
   chrome.runtime.onMessage.addListener(fn);
}
export function addTabsCreatedListener(fn) {
   chrome.tabs.onCreated.addListener(fn);
}
export function addTabsActivatedListener(fn) {
   chrome.tabs.onActivated.addListener(fn);
}
export function addTabsUpdatedListener(fn) {
   chrome.tabs.onUpdated.addListener(fn);
}
export function queryTabs(q,fn) {
   chrome.tabs.query(q,fn);
}
export function executeScript(obj,fn) {
   chrome.scripting.executeScript(obj,fn);
}
export function getTabs(tabId,fn) {
   chrome.tabs.get(tabId,fn);
}
export function sendTabsMessage(tabId,msg,opts,cb) {
   chrome.tabs.sendMessage(tabId,msg,opts,cb);
}
export function lastError() {
   return chrome.runtime.lastError
}
export function sendRuntimeMessage(tabId,msg,opts,cb) {
   chrome.runtime.sendMessage(tabId,msg,opts,cb);
}

function rejectUrl(url,group){
   if (url.match(chromeUrlsRegex)) {
      return true;
   }
   return linkSubmittedRecently(url, group);
}

function urlKey(url,group) {
   return hash(`${url}|${group}`)
}

/**
 * LinkSubmittedRecently checks to see if a URL + groupId has been submitted within the past hour.
 */
function linkSubmittedRecently(url, group) {
   let key = urlKey(url,group)
   if (!recentlySubmittedLinks.hasOwnProperty(key)) {
      return false;
   }
   let currentTime = new Date();
   let lastSubmittedTime = recentlySubmittedLinks[key];
   let milliseconds = currentTime - lastSubmittedTime;
   let hours = milliseconds / (1000 * 60 * 60);
   return hours < 1;
}

function hash(s) {
   let hash = 0, i, chr;
   if (s.length === 0) return hash;
   for (i = 0; i < s.length; i++) {
      chr = s.charCodeAt(i);
      hash = ((hash << 5) - hash) + chr;
      hash |= 0; // Convert to 32bit integer
   }
   return hash;
}

function throwNotSetError() {
   throw new Error("variable `chrome` not set. Call `setChromeObject()` to set.");
}
