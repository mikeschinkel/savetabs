document.addEventListener('DOMContentLoaded', function () {
  chrome.windows.getAll({ populate: true }, function (windows) {
    var windowMap = {};
    windows.forEach(function (window, index) {
      windowMap[window.id] = `Window${index + 1}`;
    });

    chrome.tabs.query({}, function (tabs) {
      var tabListBody = document.getElementById('tabListBody');
      var tabGroupPromises = [];
      var noTabGroupUrls = [];

      tabs.forEach(function (tab) {
        var windowName = windowMap[tab.windowId] || 'Unknown';

        var promise = new Promise(function(resolve, reject) {
          if (tab.groupId !== -1) {
            chrome.tabGroups.get(tab.groupId, function(tabGroup) {
              resolve({ windowName: windowName, tabGroup: tabGroup ? tabGroup.title : '<none>', url: tab.url });
            });
          } else {
            noTabGroupUrls.push({ windowName: windowName, url: tab.url });
            resolve(null);
          }
        });

        tabGroupPromises.push(promise);
      });

      Promise.all(tabGroupPromises).then(function(tabGroupData) {
        // Populate the table with URLs without a TabGroup
        noTabGroupUrls.forEach(function(data) {
          var row = document.createElement('tr');
          row.innerHTML = `<td>${data.windowName}</td><td><none></td><td>${data.url}</td>`;
          tabListBody.appendChild(row);
        });

        var groupedUrls = {};

        // Group URLs by TabGroup
        tabGroupData.forEach(function(data) {
          if (data) {
            if (!groupedUrls[data.tabGroup]) {
              groupedUrls[data.tabGroup] = [];
            }
            groupedUrls[data.tabGroup].push({ windowName: data.windowName, url: data.url });
          }
        });

        // Populate the table with grouped URLs
        for (var tabGroup in groupedUrls) {
          var urls = groupedUrls[tabGroup];
          urls.forEach(function(data) {
            var row = document.createElement('tr');
            row.innerHTML = `<td>${data.windowName}</td><td>${tabGroup}</td><td>${data.url}</td>`;
            tabListBody.appendChild(row);
          });
        }
      });
    });
  });
});


