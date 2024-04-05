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

export function getTabGroupResources(tabGroup) {
   return new Promise((resolve, reject) => {
      chrome.tabs.query({ groupId: tabGroup.id }, tabs => {
         if (chrome.runtime.lastError) {
            reject(chrome.runtime.lastError.message);
         } else {
            resolve(tabs
               .filter(tab =>  !rejectUrl(tab.url,tab.group) )
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

export function getAllTabGroupsResources(tabGroups) {
    // Create a promise to retrieve tabs for each tab group
    const tabGroupPromises = tabGroups.map(tabGroup => {
        return getTabGroupResources(tabGroup);
    });

    // Create a promise to retrieve tabs without a tab group
    const noGroupPromise = getTabGroupResources({
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

export async function getAllResources() {
   try {
      const tabGroups = await getAllTabGroups();
      return await getAllTabGroupsResources(tabGroups);
   } catch (error) {
      console.log('Error: getAllResources(): ', error);
      throw error;
   }
}

export function addRecentlySubmittedResources(resources) {
   let d = new Date();
   resources.forEach(r => {
      recentlySubmittedResources[urlKey(r.url,r.groupId)] = d;
   })
}

const chromeUrlsRegex = /^(chrome(-\w+)?):\/\//;
function rejectUrl(url,group){
   if (url.match(chromeUrlsRegex)) {
      return true;
   }
   return resourceSubmittedRecently(url, group);
}

function urlKey(url,group) {
   return `${url}|${group}`
}

let recentlySubmittedResources = {};
// resourceSubmittedRecently checks to see if a URL + groupId has been submitted within the past hour.
function resourceSubmittedRecently(url,group) {
   let key = urlKey(url,group)
   if (!recentlySubmittedResources.hasOwnProperty(key)) {
      return false;
   }
   let currentTime = new Date();
   let lastSubmittedTime = recentlySubmittedResources[key];
   let milliseconds = currentTime - lastSubmittedTime;
   let hours = milliseconds / (1000 * 60 * 60);
   return hours < 1;
}


