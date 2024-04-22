import {getApiServerUrl} from './api.js';
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
   Alpine.data('allCheckboxChecker', () => {
      return ({
         confirmationDialogState: 'closed',
         confirmationPrompt: '',
         checkedCheckbox: {},
         getConfirmationPrompt() {
            return this.confirmationPrompt;
         },
         showConfirmationDialog() {
            return this.confirmationDialogState === 'open';
         },
         getRowCheckboxes(obj){
            const form = obj.closest('form');
            return Array.from(form.querySelectorAll(`input[type="checkbox"]:not(.check-all)`));
         },
         getHeadOrFootCheckbox(obj){
            const tableId = obj.closest('table').id;
            const form = obj.closest('form');
            const trId = obj.closest('tr').id;
            return form.querySelectorAll(`input[type="checkbox"]:not(#${trId})`)[0];
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
            // Ask for confirmation since some rows checked, some rows not
            const action = event.target.checked ? 'select' : 'deselect';
            const state = event.target.checked ? 'unselected' : 'selected';
            this.confirmationPrompt = `Are you sure you want to ${action} all those ${state}?`;
            this.confirmationDialogState = 'open';
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
            this.confirmationDialogState = 'closed';
         },
         closeConfirmationDialog(event) {
            event.stopPropagation();
            event.preventDefault();
            this.checkedCheckbox.checked = !this.checkedCheckbox.checked
            this.confirmationDialogState = 'closed';
         },
      });
   });
});
