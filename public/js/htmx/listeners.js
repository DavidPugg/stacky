document.body.addEventListener('htmx:beforeSwap', function (event) {
  event.detail.shouldSwap = true;
});

document.body.addEventListener('updatePageDetails', function (event) {
  document.title = event.detail.title;
  document.querySelector('meta[name="description"]').content =
    event.detail.description;
});

document.body.addEventListener('redirect', function (event) {
  htmx.ajax('GET', event.detail.value);
});

document.body.addEventListener('removeNoComments', function (event) {
  const noComments = document.getElementById('no-comments');
  const commentsList = document.getElementById('comment-list');
  noComments.classList.remove('visible');
  noComments.classList.add('hidden');
  commentsList.classList.remove('hidden');
  commentsList.classList.add('visible');
});

document.body.addEventListener('addNoComments', function (event) {
  const commentsList = document.getElementById('comment-list');

  if (commentsList.childElementCount == 0) {
  const noComments = document.getElementById('no-comments');
  
  noComments.classList.remove('hidden');
  noComments.classList.add('visible');
  commentsList.classList.remove('visible');
  commentsList.classList.add('hidden');
  }
});
