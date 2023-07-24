document.addEventListener('htmx:beforeSwap', function (event) {
  event.detail.shouldSwap = true;
});
