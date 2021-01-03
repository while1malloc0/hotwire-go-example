# Hotwire Go Example

This is a recreation of the [Hotwire Rails Demo Chat](https://github.com/hotwired/hotwire-rails-demo-chat) with a Go backend.
See the [Hotwire docs](https://hotwire.dev) for more information about Hotwire.

## Quickstart

### Requirements

* Go
* NPM

### Set up the project

#### Automatic

Run `./scripts/setup.sh`.

#### Manual

1. Download Go dependencies: `go mod download`
1. Download node dependencies: `npm install`
1. Create a webpack bundle: `npx webpack -c webpack.config.js`

### Run the project

The demo can be run with `go run main.go` from the project root.

## Packages used

Ruby On Rails provides a wealth of functionality out of the box.
While batteries-included frameworks for Go web development exist (e.g. [Buffalo](https://gobuffalo.io)), 
it's often more idiomatic to compose functionality from small, single-purpose libraries. 
In order to mimic some of Rails' functionality, the following packages were used:

* [Chi](https://github.com/go-chi/chi) provides request routing and middleware support.
* [Render](https://github.com/unrolled/render) provides template rendering.
* [GORM](https://gorm.io) provides database ORM functionality.
* [Websocket](https://nhooyr.io/websocket) provides websocket connectivity.

In addition, the `pkg` directory contains a few purpose-built packages that mimic Rails functionality for the demo:

* `pubsub` implements in-memory PubSub functionality that fills the role of Rails' [ActionCable](https://edgeguides.rubyonrails.org/action_cable_overview.html).
* `notify` implements a cookie-based notification banner.
* `timefmt` implements a helper function for formatting Go's `time.Time` type into a similar format to the Rails example.

## Notable differences from the Rails demo

* The Rails demo uses a [Stimulus controller](https://stimulus.hotwire.dev) to do form resetting when submitting a new chat message.
While Stimulus is a great framework, in Go it's generally idiomatic to reimplement small pieces of functionality instead of adding a new dependency, and because the reset controller functionality is easily implemented in pure Javascript, it felt more idiomatic to leave Stimulus out in this case.
* [Turbo Rails](https://github.com/hotwired/turbo-rails), and the Rails example by extension, relies heavily on `ActionCable` to connect Turbo Streams to a websocket.
Because there is no ActionCable equivalent in Go, this example uses an in-memory pubsub mechanism and a `TurboStreamWebsocketSource` [custom element](https://developer.mozilla.org/en-US/docs/Web/Web_Components/Using_custom_elements) to mimic the same functionality.
In a production setting, you would probably want to use a more robust PubSub mechanism.