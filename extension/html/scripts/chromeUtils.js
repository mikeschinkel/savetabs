// chromeUtils.js
const notSetError = new Error("variable `chrome` not set. Call `setChromeObject()` to set.")
let chrome;
chrome = {
   windows: {
      getAll: function (tabGroupId, callback) {
         throw notSetError;
      }
   },
   tabGroups: {
      get: function (tabGroupId, callback) {
         throw notSetError;
      }
   },
   tabs: {
      query: function (queryInfo, callback) {
         throw notSetError;
      }
   },
   runtime: {
      lastError: {
         message: "chrome object not set with setChromeObject()"
      }
   }
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
                     url: tab.url,
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

export async function getAllLinks() {
   try {
      const tabGroups = await getAllTabGroups();
      return await getAllTabGroupsLinks(tabGroups);
   } catch (error) {
      console.log('Error: getAllLinks(): ', error);
      throw error;
   }
}

export function addRecentlySubmittedLinks(links) {
   let d = new Date();
   links.forEach(r => {
      recentlySubmittedLinks[urlKey(r.url,r.group)] = d;
   })
}

const chromeUrlsRegex = /^(chrome(-\w+)?):\/\//;
function rejectUrl(url,group){
   if (url.match(chromeUrlsRegex)) {
      return true;
   }
   return LinkSubmittedRecently(url, group);
}

function urlKey(url,group) {
   return hash(`${url}|${group}`)
}

let recentlySubmittedLinks = {};
// LinkSubmittedRecently checks to see if a URL + groupId has been submitted within the past hour.
function LinkSubmittedRecently(url,group) {
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

