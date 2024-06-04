
// Capture the item(s) being dragged
document.addEventListener('dragstart', function (event) {
   if (event.target === null) {
      return
   }
   const target=event.target;
   if (target.dataset === null) {
      return
   }
   const dataset = target.dataset;
   if (dataset.draggable === null) {
      return
   }
   const dt = event.dataTransfer
   dt.setData("text/plain",dataset.draggable)
   // See https://stackoverflow.com/a/62385211/102699
   dt.dropEffect = 'move';
   dt.effectAllowed = 'move';
   console.log(`${event.type}:`,dragDropItem);
});

// Call API to update DB on drop
// Also update the current links view showing those items no longer visible
document.addEventListener('drop', function (event) {
   console.log(event.type, event);
});

document.addEventListener('dragover', function (event) {
   //see https://stackoverflow.com/questions/21339924/drop-event-not-firing-in-chrome/36207613
   event.preventDefault();
});


// Display droppable element as ready to accept drag
// document.addEventListener('dragenter', function (event) {
//    console.log(event.type, event);
// });

// Unhighlight droppable element
// document.addEventListener('dragleave', function (event) {
//    console.log(event.type, event);
// });

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