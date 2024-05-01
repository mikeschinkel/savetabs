(async () => {
   await import(chrome.runtime.getURL('./html/scripts/content.js'));
})();