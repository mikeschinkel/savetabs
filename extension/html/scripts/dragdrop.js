import {apiPostOnDrop} from "./api.js";

const dragDrop = newDragDrop();

export function newDragDropEvents() {
   return {
      dragStart: function dragStart(event) {
         const target = event.target;
         dragDrop.dragElement = target;
         dragDrop.addDragElement(target);
         const item = dragDrop.validateDragElement(event.target)
         if (!item) {
            return
         }
         dragDrop.drag = item.drag;
         dragDrop.parent = item.parent;

         const dt = event.dataTransfer
         dt.setDragImage(dragDrop.getDragImage(), 10, 10);
         dt.effectAllowed = dragDrop.effectAllowed;
         dt.setData(dragDrop.mimeType(), dragDrop.value());

         const a = event.target.querySelector('a');
         if (!a) {
            return
         }
         const url = a.getAttribute('href');
         if (!url) {
            return
         }
         const title = a.getAttribute('title');
         if (title && url !== title) {
            // See https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/Recommended_drag_types#dragging_links
            dt.setData("text/x-moz-url", `${url}\n${title}`)
         }
         dt.setData("text/uri-list", url)
         dt.setData("text/plain", url)

         // See https://stackoverflow.com/a/62385211/102699
         // console.log(`${event.type}:`,dataset.dragsources);
      },
      dragOver: function dragOver(event) {
         if (!dragDrop.isDroppableTarget(event)) {
            return
         }
         //see https://stackoverflow.com/questions/21339924/drop-event-not-firing-in-chrome/36207613
         event.preventDefault();
         event.dataTransfer.dropEffect = dragDrop.effectAllowed;
         const dt = dragDrop.getDropTarget(event.target)
         dragDrop.highlightDropTarget(dt.element);
      },
      dragEnter: function dragEnter(event) {
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
         if (!dragDrop.isDroppableTarget(event)) {
            return
         }
         const dt = dragDrop.getDropTarget(et)
         if (!dt) {
            return
         }
         dragDrop.highlightDropTarget(dt.element);
      },
      dragLeave: function dragLeave(event) {
         if (!dragDrop.isDroppableTarget(event)) {
            return
         }
         if (event.target.contains(event.relatedTarget)) {
            return
         }
         dragDrop.unHighlightDropTarget(event.target)
      },
      dragEnd: function dragEnd(event) {
         dragDrop.clearDragElements()
      },
      drop: function drop(event) {
         if (!dragDrop.isDroppableTarget(event)) {
            return
         }
         const et = event.target
         const drop = dragDrop.validateDrop(et);
         if (!drop) {
            return
         }
         dragDrop.drop = drop;

         let elements = dragDrop.captureElements();
         apiPostOnDrop(dragDrop, function (status, response) {
            if (!status) {
               alert(response); // TODO: Create a better UX for this
               return false;
            }
            dragDrop.removeCapturedElements(elements);
            dragDrop.clearDragElements();
            return true;
         })
         dragDrop.unHighlightDropTarget(et)
      }
   }
}

function newDragDrop() {
   const parent = {
      type: '',
      id: 0
   }
   const drag = {
      type: '',
      ids: [],
   }
   const drop = {
      type: '',
      id: 0
   }
   return {
      mimeTypePrefix: 'application/x-savetabs',
      imageSrc: "./assets/savetabs-32.png",
      effectAllowed: 'move',
      dropTargetWrapper: 'details',
      dragParentAttr: 'data-dragparent',
      dragSourcesAttr: 'data-dragsources',
      dropTypeAttr: 'data-droptypes',
      dropTargetAttr: 'data-droptarget',
      highlightedDropTargets: [],
      onHoverClasses: [
         'backdrop-invert',
         'backdrop-opacity-20'
      ],
      dragElement: null,
      dragElements: [],
      parent: parent,
      drag: drag,
      drop: drop,
      addDragElement(element) {
         dragDrop.dragElements.push(element);
      },
      mimeType() {
         return this.getMimeType(`${this.parent.type}+${this.drag.type}`);
      },
      params() {
         const params = {
            parent_type: encodeURIComponent(this.parent.type),
            parent_id: parseInt(this.parent.id),
            drag_type: encodeURIComponent(this.drag.type),
            drag_id: this.parseInts(this.drag.ids),
            drop_type: encodeURIComponent(this.drop.type),
            drop_id: parseInt(this.drop.id, 10)
         }
         return new URLSearchParams(params);
      },
      getDragImage: function () {
         const img = new Image();
         img.src = this.imageSrc;
         return img;
      },
      value: function () {
         return this.params().toString();
      },
      getMimeType(ddType) {
         return `${this.mimeTypePrefix}-${ddType}`;
      },
      getPostData() {
         return {
            parent: {
               type: encodeURIComponent(this.parent.type),
               id: parseInt(this.parent.id, 10),
            },
            drag: {
               type: encodeURIComponent(this.drag.type),
               ids: this.parseInts(this.drag.ids),
            },
            drop: {
               type: encodeURIComponent(this.drop.type),
               id: parseInt(this.drop.id, 10)
            }
         }
      },
      getDropTarget(target) {
         let dropTypes = null;
         while (!dropTypes) {
            dropTypes = target.getAttribute(this.dropTypeAttr)
            if (dropTypes) { // TODO: This check will not be sufficient when more than one are supported.
               break
            }
            target = target.parentElement
            if (!target) {
               return null
            }
            if (target.localName === this.dropTargetWrapper) {
               return null
            }
         }
         const value = target.getAttribute(this.dropTargetAttr)
         if (!value) {
            console.log(`Missing or incorrectly formatted '${this.dropTargetAttr}' for`, event.target)
            return null
         }
         const [dropType, dropId] = value.split(':');
         return {
            types: dropTypes,
            element: target,
            value: value,
            type: dropType,
            id: parseInt(dropId,10)
         }
      },
      parseInts(nums) {
         let ints = []
         for (let num of nums) {
            let i = parseInt(num, 10)
            if (i === 0 && num !== "0") {
               console.log(`Invalid value '${num}' parsed as int`)
               continue
            }
            ints.push(i)
         }
         return ints;
      },
      validateDrop(target) {
         const dt = this.getDropTarget(target)
         if (!dt) {
            return false;
         }
         return {
            type: dt.type,
            id: parseInt(dt.id,10)
         };
      },
      validateDragElement(target) {
         const dragSources = target.getAttribute(this.dragSourcesAttr)
         if (!dragSources) {
            console.log(`Attribute '${this.dragSourcesAttr}' missing from ${target.id}`)
            return false;
         }
         const dragParent = target.getAttribute(this.dragParentAttr)
         if (!dragParent) {
            console.log(`Attribute '${this.dragParentAttr}' missing from ${target.id}`)
            return false;
         }
         const [dragType, idOfType] = dragSources.split(":")
         if (!idOfType) {
            console.log(`Attribute '${this.dragSourcesAttr}'='${dragSources}' in wrong format; should be '<type>:<id>' where <type> can be 'link', 'group', etc.`)
            return false;
         }
         const [parentType, idOfParent] = dragParent.split(":")
         if (!idOfParent) {
            console.log(`Attribute '${this.dragParentAttr}'='${dragParent}' in wrong format; should be '<type>:<id>' where <type> can be 'group', etc.`)
            return false;
         }
         return {
            drag: {
               type: dragType,
               ids: this.parseInts([idOfType])
            },
            parent: {
               type: parentType,
               id: parseInt(idOfParent, 10)
            }
         };
      },
      captureElements() {
         // Capture DOM elements to allow us to delete them after
         // `dragend` even removes them from `dragElements` array.
         const capture = [];
         for (let e of dragDrop.dragElements) {
            capture.push({
               parent: e.parentNode,
               child: e
            })
         }
         return capture;
      },
      removeCapturedElements(capture) {
         for (let e of capture) {
            if (!e) {
               alert(`Problem with capture: ${JSON.stringify(capture)}`)
            }
            e.parent.removeChild(e.child);
         }
      },
      clearDragElements() {
         this.dragElements = [];
      },
      highlightDropTarget(element) {
         if (!element) {
            return;
         }
         this.unHighlightDropTarget();
         //console.log(`Adding highlight: ${element.localName}`);
         for (let cls of this.onHoverClasses){
            element.classList.add(cls);
         }
         this.highlightedDropTargets.push(element);
      },
      unHighlightDropTarget(element) {
         if (element) {
            this.highlightedDropTargets.push(element);
         }
         for (let dt of this.highlightedDropTargets) {
            if (!dt){
               continue;
            }
            //console.log(`Removing highlight: ${dt.localName}`);
            for (let cls of this.onHoverClasses){
               dt.classList.remove(cls);
            }
         }
         this.highlightedDropTargets = [];
      },
      isDroppableTarget(event) {
         let dt = this.getDropTarget(event.target);
         if (!dt) {
            return null;
         }
         if (!(dt.types)) {
            return null;
         }
         if (dt.id===this.parent.id) {
            return null;
         }
         const dts = dt.types.split(' ');
         const includes = (dt) => event.dataTransfer.types.includes(dt);
         for (const _dt of dts) {
            let result = includes(this.getMimeType(_dt));
            if (!result) {
               // Try the raw type if the SaveTab's specific type did not work
               result = includes(_dt);
            }
            if (!result) {
               continue
            }
            return result
         }
         return null
      },
   }
}
