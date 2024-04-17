
// Load Alpine.js after initial configurations are done
function loadAlpine(caller) {
   const script = document.createElement('script');
   script.src = 'scripts/alpinejs-csp@3.13.8.min.js';
   script.integrity = 'sha384-ngWUWXGZwWaMILMRRZR1rmY/dFdl79tJTC1pJD+w4Ca0uyMJOf0p+L5C1fkWht0l';
   script.crossorigin = 'anonymous';
   script.defer = true;
   script.onload = () => {
      console.log('Alpine.js has loaded, initializing now...');
   };
   document.head.appendChild(script);
}

window.onload = function() {
   loadAlpine();  // This ensures everything is fully loaded
};
