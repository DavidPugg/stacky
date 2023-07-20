document.body.addEventListener('htmx:beforeSwap', function (event) {
  if (event.detail.isError) {
    event.detail.shouldSwap = true;
    event.detail.target = document.body;
    event.detail.swapStyle = 'innerHTML';
  }
});
