import {getApiServerUrl} from './api.js';
import {} from './shared.js'
import {} from './menu.js'
import {} from './icons.js'
import {} from './alpine-loader.js'

console.log("SaveTabs daemon:", getApiServerUrl())

document.addEventListener('alpine:init', () => {
   Alpine.data('checkedHighlighter', () => ({
      highlight(event) {
         const input = event.target;
         const tr = input.closest('tr');
         const highlight = tr.dataset.highlight;
         const classes = tr.classList;
         if (input.checked) {
            classes.add(highlight);
         } else {
            classes.remove(highlight);
         }
      }
   }));
   Alpine.data('checkboxChecker', () => {
      return ({
         confirmDialogState: 'closed',
         confirmPrompt: '',
         checkedCheckbox: {},
         getConfirmPrompt() {
            return this.confirmPrompt;
         },
         showConfirmDialog() {
            return this.confirmDialogState === 'open'
         },
         getRowCheckboxes(obj){
            const form = obj.closest('form');
            return Array.from(form.querySelectorAll(`input[type="checkbox"]:not(.check-all)`));
         },
         getHeadOrFootCheckbox(obj){
            const tableId = obj.closest('table').id;
            const form = obj.closest('form');
            const trId = obj.closest('tr').id;
            return form.querySelectorAll(`input[type="checkbox"].check-all:not(tr#${trId} input)`)[0];
         },
         allChecked(checkboxes) {
            return checkboxes.every(_ =>  _.checked);
         },
         noneChecked(checkboxes) {
            return checkboxes.every(_ =>  !_.checked);
         },
         maybeConfirmCheckAll(event) {
            const checkboxes = this.getRowCheckboxes(event.target)
            event.stopPropagation();
            event.preventDefault();
            this.checkedCheckbox =  event.target;
            if (this.checkedCheckbox.checked && this.noneChecked(checkboxes)) {
               this.changeAllCheckboxes(event);
               return
            }
            if (!this.checkedCheckbox.checked && this.allChecked(checkboxes)) {
               this.changeAllCheckboxes(event);
               return
            }
            const classList = this.$refs.confirmDialog.classList
            classList.add("modal-open")
            classList.add("modal")
            classList.remove("hidden")
            // Ask for confirm since some rows checked, some rows not
            const action = event.target.checked ? 'select' : 'deselect';
            this.confirmPrompt = `Are you sure you want to ${action} ALL?`;
            this.confirmDialogState = 'open';
            return true;
         },
         changeAllCheckboxes(event) {
            const checkbox = this.checkedCheckbox;
            event.stopPropagation();
            event.preventDefault();
            this.getHeadOrFootCheckbox(checkbox).checked = checkbox.checked;
            this.getRowCheckboxes(checkbox).forEach(function (_) {
               _.checked = checkbox.checked;
               _.dispatchEvent(new Event('click'));
            });
            this.setClosed();
         },
         setClosed() {
            const classList = this.$refs.confirmDialog.classList
            classList.add("hidden")
            classList.remove("modal")
            classList.remove("modal-open")
            this.confirmDialogState = 'closed';
         },
         closeConfirmDialog(event) {
            event.stopPropagation();
            event.preventDefault();
            this.checkedCheckbox.checked = !this.checkedCheckbox.checked
            this.setClosed();
         },
      });
   });
});

/**
 * Capture data-action value on `<input type="submit"> and set the value of `<input name="action">` before HTMX request.
 */
document.addEventListener('htmx:trigger', function(event) {
   const el = document.activeElement;
   const {tagName, type} = el;
   if (tagName !== 'INPUT') {
      return
   }
   if (type !== 'submit') {
      return
   }
   const form = el.closest('form')
   if (!form) {
      alert(`<form> not found for ${el.name}`);
      return
   }
   const input = form.querySelector('input[type="hidden"][name="action"]');
   if (!input) {
      alert(`<input name="action"> not found for ${el.name}`);
      return
   }
   const action = el.getAttribute('data-action')
   if (action==="") {
      alert(`<input name="action" data-action="..."> not found for ${el.name}`);
      return
   }
   input.value = action;
});

const htmxEvents = [
   // 'trigger',
   // 'confirm',
   // 'validate',
   // 'configRequest',
   // 'validateUrl',
   // 'beforeRequest',
   // 'beforeSend',
   // 'xhr:loadstart',
   // 'xhr:progress',
   // 'xhr:loadend',
]
// for (let index in htmxEvents) {
//    let ev = htmxEvents[index];
//    document.body.addEventListener(`htmx:${ev}`, function (event) {
//       console.log(event);
//    });
// }
