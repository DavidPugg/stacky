document.addEventListener('htmx:responseError', function (evt) {
  if (evt.detail.requestConfig.headers['Replace-Content'] == 'true') {
    document.documentElement.innerHTML = evt.detail.xhr.response;
  }
});
