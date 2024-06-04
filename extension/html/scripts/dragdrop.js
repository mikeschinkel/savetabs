const mimeTypePrefix = 'application/x-savetabs';

function dragDropMimeType(dragType) {
   return `${mimeTypePrefix}-${dragType}`
}
// Capture the item(s) being dragged
document.addEventListener('dragstart', function (event) {
   const target=event.target;
   if (target.dataset === null) { // TODO: Verify that checking null is correct here
      console.log(`Attributes data-draggable missing from ${event.target.id}`)
      return
   }
   const dataset = target.dataset;
   if (dataset.draggable === null) { // TODO: Verify that checking null is correct here
      console.log(`Attribute data-draggable missing from ${event.target.id}`)
      return
   }
   const [dragType,idOfType] = dataset.draggable.split(":")
   if (!idOfType) {
      console.log(`Attribute data-draggable='${dataset.draggable}' in wrong format; should be '<type>:<id>' where <type> can be 'link', 'group', etc.`)
      return
   }
   const mimeType = dragDropMimeType(dragType);
   const dt = event.dataTransfer
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
   console.log(`${event.type}:`,dataset.draggable);
});

document.addEventListener('dragover', function (event) {
   if (!hasDroppedItem(event)) {
      return
   }
   //see https://stackoverflow.com/questions/21339924/drop-event-not-firing-in-chrome/36207613
   event.preventDefault();
   const et = event.target;
   et.dropEffect = 'move';
   console.log(`Dragover ${et.localName}#${et.id}.${et.className}`);
});

// Display droppable element as ready to accept drag
document.addEventListener('dragenter', function (event) {
   console.log(event.type, event.target.localName);
   if (!hasDroppedItem(event)) {
      return
   }
   const cl = event.target.classList
   cl.add('border-black');
   cl.add('border-4');
   console.log(event.type, event);
});
function unHighlightDroppable(event) {
   const cl = event.target.classList
   cl.remove('border-black');
   cl.remove('border-4');
}
// Unhighlight droppable element
document.addEventListener('dragleave', function (event) {
   if (!hasDroppedItem(event)) {
      return
   }
   unHighlightDroppable(event)
   console.log(event.type, event);
});

// Call API to update DB on drop
// Also update the current links view showing those items no longer visible
document.addEventListener('drop', function (event) {
   const item = getDroppedItem(event);
   if (!item) {
      console.log("Unexpected empty item in 'drop' event", event)
      return
   }
   const [dragType,dragId] = item.split(':');
   console.log("Drop:", `${dragType}:${dragId}`);
   unHighlightDroppable(event)
});


function hasDroppedItem(event){
   return _processDroppedItem(event, function (event,dragType) {
      return event.dataTransfer.types.includes(dragType)
   });
}

function getDroppedItem(event){
   return _processDroppedItem(event, function (event,dragType) {
      let item = event.dataTransfer.getData(dragType)
      if (item) {
         item = `${dragType}:${item}`;
      }
      return item;
   });
}
function _processDroppedItem(event, func){
   const et = event.target
   const droppable = et.getAttribute('data-droppable');
   if (!droppable) {
      return null
   }
   const dragTypes = droppable.split(' ');
   const dt = event.dataTransfer
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