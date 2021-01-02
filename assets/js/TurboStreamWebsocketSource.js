import { connectStreamSource, disconnectStreamSource } from "@hotwired/turbo";

export default class TurboStreamWebsocketSource extends HTMLElement {
  get src() {
    return this.getAttribute("src");
  }

  set src(value) {
    if (value) {
      this.setAttribute("src", value);
    } else {
      this.removeAttribute("src");
    }
  }

  async connectedCallback() {
    connectStreamSource(this);
    this.ws = this.setupWebsocket();
  }

  disconnectedCallback() {
    disconnectStreamSource(this);
    if (this.ws) {
      this.ws = null;
    }
  }

  dispatchMessageEvent(messageEvent) {
    const event = new MessageEvent("message", { data: messageEvent.data });
    this.dispatchEvent(event);
  }

  setupWebsocket() {
    let socketLocation = `ws://${window.location.host}${this.src}`
    let ws = new WebSocket(socketLocation);
    ws.onmessage = (msg) => {
      this.dispatchMessageEvent(msg);
    };
    return ws;
  }
}
