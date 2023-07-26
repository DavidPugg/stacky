document.body.addEventListener('htmx:beforeSwap', function (event) {
  event.detail.shouldSwap = true;
});

document.body.addEventListener('showAlert', function (event) {
  htmx.ajax(
    'GET',
    `/utils/showAlert?message=${event.detail.message}&type=${event.detail.type}`,
  );
});

document.body.addEventListener('updatePageDetails', function (event) {
  document.title = event.detail.title;
  document.querySelector('meta[name="description"]').content =
    event.detail.description;
});
