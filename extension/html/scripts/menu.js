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
   Alpine.data('contextMenuer', () => ({
      show(event) {
         const hidden = "hidden"
         const cm = this.getContextMenu()
         const style = cm.style;
         event.preventDefault()
         style.left = `${event.pageX}px`;
         style.top = `${event.pageY}px`;
         cm.classList.remove(hidden);
      },
      getContextMenu() {
         return document.getElementById('context-menu');
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
