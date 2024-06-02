let clicked = false;

window.isBranchCollapsed = (id) => {
   const el = document.getElementById(id)
   if (el === null) {
      return true
   }
   const els = el.getElementsByTagName('ul')
   if (els.length === 0) {
      return true
   }
   const display = els[0].style.display
   if (display === 'none') {
      return true
   }
   if (display === '') {
      return true
   }
   return false;
}

document.addEventListener('alpine:init', () => {
   const hidden = "hidden"
   Alpine.data('contextMenuer', () => ({
      contextMenu: document.getElementById('context-menu'),
      targetItem: null,
      originalValue: '',
      init() {
         this.contextMenu.addEventListener('close',this.onClose.bind(this))
      },
      focusInput() {
         this.targetItem.removeEventListener('htmx:afterSettle', this.focusInput)
         if (this.targetItem === null) {
            return
         }
         const input = this.targetItem.querySelector('input');
         if (input==null) {
            return
         }
         this.originalValue = input.value
         this.$nextTick(() => input.select())
      },
      submit(event) {
      },
      show(event) {
         event.preventDefault()

         this.targetItem = event.currentTarget
         this.targetItem.addEventListener('htmx:afterSettle', this.focusInput.bind(this))

         const style = this.contextMenu.style
         style.marginLeft = `${event.pageX}px`;
         style.marginTop = `${event.pageY}px`;

         const classList = event.currentTarget.classList
         classList.add('bg-gray-400')
         classList.add('text-white')
         this.contextMenu.showModal();
      },
      onClose(event) {
         if (this.targetItem !== null) {
            const classList = this.targetItem.classList
            classList.remove('bg-gray-400')
            classList.remove('text-white')
         }
      },
      hide(event) {
         this.contextMenu.close();
      }
   }));
   Alpine.data('preventable', () => ({
      preventExpandOnIconClick: function (event) {
         if (['svg','path'].includes(event.target.tagName.toLowerCase())) {
            return;
         }
         event.preventDefault()
      }
   }))
   Alpine.data('collapsible', () => ({
      state: 'collapsed',
      expanded: function () {
         return this.state === 'expanded'
      },
      collapsed: function () {
         return this.state === 'collapsed'
      },
      toggle: function () {
         switch (this.state) {
            case 'collapsed':
               this.state = 'expanded';
               break;
            case 'expanded':
               this.state = 'collapsed';
               break;
            default:
               alert(`Unexpected collapsible.state: ${this.state}`);
         }
      },
   }))
})


