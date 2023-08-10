class App {
  static navigate(path) {
    htmx.ajax('GET', path);
  }
}

window.App = App;
