const mimeTypePrefix = 'application/x-savetabs';
const dropTargetWrapper = 'details';
const img = new Image();
img.src = "./assets/savetabs-32.png";

// Capture the item(s) being dragged
document.addEventListener('dragstart', function (event) {
   const target=event.target;
   if (target.dataset === undefined) { // TODO: Verify that checking null is correct here
      console.log(`Attributes data-dragsources missing from ${event.target.id}`)
      return
   }
   const dataset = target.dataset;
   if (dataset.dragsources === undefined) { // TODO: Verify that checking null is correct here
      console.log(`Attribute data-dragsources missing from ${event.target.id}`)
      return
   }
   const [dragType,idOfType] = dataset.dragsources.split(":")
   if (!idOfType) {
      console.log(`Attribute data-dragsources='${dataset.dragsources}' in wrong format; should be '<type>:<id>' where <type> can be 'link', 'group', etc.`)
      return
   }
   const mimeType = dragDropMimeType(dragType);
   const dt = event.dataTransfer
   dt.setDragImage(img, 10, 10);
   dt.effectAllowed = 'move';
   dt.setData(mimeType,idOfType);

   const a = event.target.querySelector('a');
   if (!a) {
      return
   }
   const url = a.getAttribute('href');
   if (!url) {
      return
   }
   const title = a.getAttribute('title');
   if (title) {
      // See https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/Recommended_drag_types#dragging_links
      dt.setData("text/x-moz-url",`${url}\n${title}`)
   }
   dt.setData("text/uri-list",url)
   dt.setData("text/plain",url)

   // See https://stackoverflow.com/a/62385211/102699
   // console.log(`${event.type}:`,dataset.dragsources);
});

document.addEventListener('dragover', function (event) {
   if (!hasDroppedItem(event)) {
      return
   }
   //see https://stackoverflow.com/questions/21339924/drop-event-not-firing-in-chrome/36207613
   event.preventDefault();
   event.dataTransfer.dropEffect = 'move';
   const et = event.target;
   // console.log(`Dragover ${et.localName}#${et.id}.${et.className}`);
   highlightDroppable(event.target);
});

// Display droppable element as ready to accept drag
document.addEventListener('dragenter', function (event) {
   const ert = event.relatedTarget;
   if (!ert) {
      return
   }
   const et = event.target;
   if (et.localName === ert.localName) {
      // Chances are very high that if they target and relatedTarget
      // have the same element type then it is not a drop target.
      return
   }
   if (!hasDroppedItem(event)) {
      return
   }
   // console.log(event.type, event);
   const droppable = getDroppableElement(et)
   if (!droppable) {
      return
   }
   highlightDroppable(droppable);
});


// Unhighlight droppable element
document.addEventListener('dragleave', function (event) {
   if (!hasDroppedItem(event)) {
      return
   }
   if (event.target.contains(event.relatedTarget)) {
      return
   }
   unHighlightDroppable(event.target)
   // console.log(event.type, event);
});

// Call API to update DB on drop
// Also update the current links view showing those items no longer visible
document.addEventListener('drop', function (event) {
   const source = getDroppedItem(event);
   if (!source) {
      console.log("Unexpected empty item in 'drop' event", event.target.id)
      return
   }
   const dropItem = getDroppableElement(event.target);
   if (!dropItem) {
      console.log("Missing or incorrectly formatted 'data-droptypes' for", event.target.id)
      return
   }
   const target = dropItem.getAttribute('data-droptarget')
   if (!target) {
      console.log("Missing or incorrectly formatted 'data-droptarget' for", event.target.id)
      return
   }
   let [dragType,dragId] = source.split(':');
   dragType = stripPrefix(dragType,`${mimeTypePrefix}-`);
   const [dropType,dropId] = target.split(':');
   console.log("Drag and Drop:", `${dragType}:${dragId} ==> ${dropType}:${dropId}`);
   unHighlightDroppable(event.target)
});

function stripPrefix(str, prefix) {
   if (str.startsWith(prefix)) {
      return str.slice(prefix.length);
   }
   return str;
}

function hasDroppedItem(event){
   return _processDroppedItem(event.target, function (event,dragType) {
      return event.dataTransfer.types.includes(dragType)
   });
}

function getDroppedItem(event){
   return _processDroppedItem(event.target, function (event,dragType) {
      let item = event.dataTransfer.getData(dragType)
      if (item) {
         item = `${dragType}:${item}`;
      }
      return item;
   });
}
function _processDroppedItem(target, func){
   let droppable = getDroppableType(target);
   if (!droppable) {
      return null;
   }
   const dragTypes = droppable.split(' ');
   for (const dragType of dragTypes) {
      let result = func(event,dragDropMimeType(dragType));
      if (!result) {
         // Try the raw type if the SaveTab's specific type did not work
         result = func(event,dragType);
      }
      if (!result) {
         continue
      }
      return result
   }
   return null
}

function getDroppableType(target){
   return _getDroppable('t',target)
}

function getDroppableElement(target){
   return _getDroppable('e',target)
}

function _getDroppable(result,target){
   let dropTypes = null;
   while (!dropTypes) {
      dropTypes = target.getAttribute('data-droptypes')
      if (dropTypes) {
         break
      }
      target = target.parentElement
      if (!target) {
         return null
      }
      if (target.localName===dropTargetWrapper) {
         return null
      }
   }
   switch (result) {
      case 'e': // e=element
         result = target;
         break;
      case 't': // t=type
         result = dropTypes;
         break;
   }
   return result
}

function highlightDroppable(target) {
   const dropItem = getDroppableElement(target);
   if (!dropItem){
      return
   }
   const cl = dropItem.classList
   cl.add('border-black');
   cl.add('border-4');
}

function unHighlightDroppable(target) {
   const dropItem = getDroppableElement(target);
   if (!dropItem){
      return
   }
   const cl = dropItem.classList
   cl.remove('border-black');
   cl.remove('border-4');
}

function dragDropMimeType(dragType) {
   return `${mimeTypePrefix}-${dragType}`
}

// ???
// document.addEventListener('dragend', function (event) {
//    console.log(event.type, event);
// });
//
//

// const dragEvents = [
//    'drag',
//    'dragend',
//    'dragenter',
//    'dragleave',
//    'dragover',
//    'dragstart',
//    'drop'
// ]
// for (let index in dragEvents) {
//    let ev = dragEvents[index];
//    document.addEventListener(`${ev}`, function (event) {
//       console.log(event.type, event);
//    });
// }