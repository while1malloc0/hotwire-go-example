import { connectStreamSource, disconnectStreamSource } from "@hotwired/turbo";

/**
 * Creates a persistent connection to a websocket for use with Turbo Streams
 */
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

  /**
   * Called when the element is inserted into the DOM.
   * Connects to Turbo Streams as a source and sets up the websocket connection
   * for streaming updates to turbo streams
   */
  async connectedCallback() {
    connectStreamSource(this);
    this.ws = this.setupWebsocket();
  }

  /**
   * Called when the element is removed from the DOM.
   * Disconnects from Turbo Streams and deletes the WebSocket
   */
  disconnectedCallback() {
    disconnectStreamSource(this);
    if (this.ws) {
      this.ws = null;
    }
  }

  /**
   * Called in response to a websocket message. Unpacks the websocket message
   * and dispatches it as a new MessageEvent to Turbo Streams.
   * 
   * @param {MessageEvent} messageEvent The original message to dispatch
   */
  dispatchMessageEvent(messageEvent) {
    const event = new MessageEvent("message", { data: messageEvent.data });
    this.dispatchEvent(event);
  }

  /**
   * Creates a WebSocket by combining the window host and the element's src attribute and wires it to dispatch messages
   */
  setupWebsocket() {
    let socketLocation = `ws://${window.location.host}${this.src}`
    let ws = new WebSocket(socketLocation);
    ws.onmessage = (msg) => {
      this.dispatchMessageEvent(msg);
    };
    return ws;
  }
}
