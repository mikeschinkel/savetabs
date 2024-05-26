import {apiPutLinkByTab, apiPutLabel} from './api.js';
import {setChromeObject, lastError, sendRuntimeMessage} from './chromeUtils.js';

console.log("SaveTabs content.js loaded.");
setChromeObject(chrome)

console.log("Chrome:", chrome);

function getContentMenu() {
   return document.getElementById('context-menu');
}

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

document.addEventListener('contextmenu', function (event) {
   const el = event.target;

   if (el.matches('[id$="-label"]')) {
      event.preventDefault();

      const id = el.id.replace('-label', '');
      const contextMenu = getContentMenu();

      contextMenu.style.left = `${event.pageX}px`;
      contextMenu.style.top = `${event.pageY}px`;
      contextMenu.classList.remove('hidden');

      // Store the target ID in the context menu for later use
      contextMenu.dataset.targetId = id;
   }
});

document.addEventListener('click', function () {
   const classList = getContentMenu().classList;
   if (!classList.contains('hidden')) {
      classList.add('hidden');
   }
});

getContentMenu().addEventListener('click', function () {
   const id = this.dataset.targetId;
   const label = document.getElementById(`${id}-label`);
   const input = document.getElementById('edit-label');

   input.value = label.textContent;
   input.style.left = `${label.offsetLeft}px`;
   input.style.top = `${label.offsetTop}px`;
   input.style.width = `${label.offsetWidth}px`;
   input.classList.remove('hidden');
   input.focus();

   // Hide the context menu
   this.classList.add('hidden');

  // Save changes on blur or Enter key, and cancel on Esc key
   input.addEventListener('blur', saveChanges);
   input.addEventListener('keydown', handleKeydown);

  function handleKeydown(event) {
    switch (event.key) {
      case 'Enter':
        saveChanges();
        break;
      case 'Escape':
        cancelChanges();
        break;
      default:
        break;
    }
  }

  function saveChanges() {
    const newText = input.value;
    label.textContent = newText;
   apiPutLabel(label)
    // Send the updated text to the server
    htmx.ajax('PUT', '/update-label', {
      target: 'body',
      values: { id: id, text: newText }
    });

    cleanup();
  }

  function cancelChanges() {
    cleanup();
  }

  function cleanup() {
    input.classList.add('hidden');
    input.removeEventListener('keydown', handleKeydown);
    input.removeEventListener('blur', saveChanges);
    input.removeEventListener('keydown', cancelChanges);
  }
});
