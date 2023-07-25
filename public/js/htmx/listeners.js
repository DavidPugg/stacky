document.body.addEventListener('htmx:beforeSwap', function (event) {
  event.detail.shouldSwap = true;
});

document.body.addEventListener('showAlert', function (event) {
  htmx.ajax(
    'GET',
    `/showAlert?message=${event.detail.message}&type=${event.detail.type}`,
  );
});
