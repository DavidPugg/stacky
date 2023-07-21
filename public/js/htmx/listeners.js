document.body.addEventListener('htmx:responseError', function (evt) {
  if (evt.detail.xhr.status === 404) {
    window.location.href = '/404';
  }

  if (evt.detail.xhr.status === 500) {
    window.location.href = '/500';
  }
});
