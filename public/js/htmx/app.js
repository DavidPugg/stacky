class App {
  static navigate(path) {
    htmx.ajax('GET', path);
  }

  static handlePostLike(liked, postId)  {
    htmx.ajax(liked == "true" ? "DELETE" : "POST",`${window.location.origin}/posts/${postId}/like`, { target: `#like-button-${postId}`, swap: "outerHTML" })
  }
}

window.App = App;
