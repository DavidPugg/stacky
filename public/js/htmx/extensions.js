(function () {
  htmx.defineExtension("reset-on-success", {
    onEvent: function (name, event) {
      if (name !== "htmx:beforeSwap") return;
      if (event.detail.isError) return;

      const triggeringElt = event.detail.requestConfig.elt;
      if (
        !triggeringElt.closest("[hx-reset-on-success]") &&
        !triggeringElt.closest("[data-hx-reset-on-success]")
      )
        return;

      if (triggeringElt.tagName === "FORM") {
        triggeringElt.reset();
      }
    },
  });
})();
