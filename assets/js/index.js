import TurboStreamWebsocketSource from "./TurboStreamWebsocketSource";

window.onload = () => {
  // This was implemented using a Stimulus controller in the original Rails
  // example. Clears the new message form on send.
  addEventListener("turbo:submit-end", () => {
    document.getElementById("new-message-form").reset();
  });

  customElements.define(
    "turbo-stream-websocket-source",
    TurboStreamWebsocketSource
  );
};
