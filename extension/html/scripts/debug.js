htmx.config.historyEnabled = true;
htmx.config.debug = true;
document.addEventListener('htmx:afterRequest', function(evt) {
   console.log("HTMX request completed:", evt.detail);
});
document.addEventListener('htmx:sendError', function(evt) {
   console.error("HTMX send error:", evt.detail);
});
