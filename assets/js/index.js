import TurboStreamWebsocketSource from "./TurboStreamWebsocketSource";

window.onload = () => {
  addEventListener("turbo:submit-end", () => {
    document.getElementById("new-message-form").reset();
  });

  customElements.define(
    "turbo-stream-websocket-source",
    TurboStreamWebsocketSource
  );
};
