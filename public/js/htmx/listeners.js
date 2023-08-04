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

document.body.addEventListener('setLoggedInUser', function () {
  htmx.ajax('POST', `/auth/set_user`);
});

document.body.addEventListener('redirect', function (event) {
  htmx.ajax('GET', event.detail.value);
});

document.body.addEventListener('removeNoComments', function (event) {
  const noComments = document.getElementById('no-comments');
  if (noComments) {
    noComments.remove();
  }
});

