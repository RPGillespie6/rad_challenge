# RAD Chat

This repo contains a simple poll-based solution to the RAD Formation coding challenge.

# Quick Start

Simply clone this repository and run `go run main.go` in the root directory of the repo. Then navigate to `localhost:3000` in one or more browser tabs.

# Design Assumptions

These are the assumptions I made when building this solution:

- Must use golang for the backend

For the backend I used the default version of go provided by Ubuntu 20.04 (Linux) repositories. In an effort to become more familiar with "vanilla go" I constrained myself to only system libraries. There are likely 3rd party libraries on GitHub that may have simplified things, but I chose to use only the stock `net/http` and `encoding/json` packages.

- Must use a poll-based syncing solution in the front end (as opposed to WebSockets)

WebSockets are the most obvious choice for a true real-time chat app because of the ability for the server to send asynchronous messages to the client without the need for the client to poll. This requires 3rd party libraries for both the front end and backend, and is also slightly more difficult to debug (especially in IE11), so I avoided this in favor of the stated requirement.

- Must work in Internet Explorer 11

You mentioned "customer facing" web applications must be compatible with IE11, so for that reason I chose to use jQuery in the front end so I could have access to a basic set of battle-tested APIs that work in all browsers. If IE11 compatibility was not a requirement I would probably forego it in favor of [fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) and other APIs provided by modern browsers (but not IE11). Note IE11 *does* have WebSocket support, surprisingly, but I didn't use WebSockets for the reasons outlined in the previous point.

- No external JS libraries

With the exception of jQuery for reasons outlined in the previous point, I avoided using external JS libraries.

- CSS/Icon libraries OK

CSS is definitely one of my weaker areas. I could have used a small .gif image for use as the upvote icon and used flexbox to lay out everything in order to avoid 3rd party CSS/Icon libraries. But for ease of use I decided to pull in Bulma and FontAwesome just so I could quickly prototype this app without having to create graphics/tinker with css.