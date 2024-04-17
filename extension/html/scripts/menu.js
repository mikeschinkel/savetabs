let clicked = false;

window.isBranchCollapsed = (id) => {
   let el = document.getElementById(id)
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
