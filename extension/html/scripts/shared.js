document.addEventListener('alpine:init', () => {
   Alpine.data('dismissibleAlert', () => ({
      dismiss(event) {
         event.target.closest('.alert').remove()
      }
   }));
});