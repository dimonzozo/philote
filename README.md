# Philote:  plug-and-play websockets server ![Build status](https://travis-ci.org/dimonzozo/philote.svg)

Philote is a minimal solution to the websockets server problem, it implements Publish/Subscribe and has a simple authentication mechanism that accomodates browser clients securely as well as server-side or desktop applications.

Simplicity is one of the design goals for Philote, ease of deployment is another: you should be able to drop the binary in any internet-accessible server and have it operational.

## Basics

Philote implements a basic topic-based [Publish-subscribe pattern](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern), messages sent over the websocket connection are classified into `channels`, and each connection is given read/write access to a given list of channels at authentication time.

Messages sent over a connection for a given channel (to which it has write permission) will be received by all other connections (that have read permission to the channel in question).

### Configuration options

Philote takes configuration options from your environment and attempts to provide sensible defaults, these are the environment variables you can set to change its behaviour:

| Environment Variable    | Default                   | Description                                                                                                        |
|:-----------------------:|:-------------------------:|:-------------------------------------------------------------------------------------------------------------------|
| `SECRET`                | ` `                       | Secret salt used to sign authentication tokens                                                                     |
| `PORT`                  | `6380`                    | Port in which to serve websocket connections                                                                       |
| `LOGLEVEL`              | `info`                    | Verbosity of log output, valid options are [debug,info,warning,error,fatal,panic]                                  |
| `MAX_CONNECTIONS`       | `255`                     | Maximum amount of concurrent websocket connections allowed                                                         |
| `READ_BUFFER_SIZE`      | `1024`                    | Size of the websocket read buffer, for most cases the default should be okay.                                      |
| `WRITE_BUFFER_SIZE`     | `1024`                    | Size of the websocket write buffer, for most cases the default should be okay.                                     |
| `CHECK_ORIGIN`          | `false`                   | Check Origin headers during WebSocket upgrade handshake.                                                           |

## Clients

* [JavaScript (browser)](https://github.com/pote/philote-js)
* [Go](https://github.com/pote/philote-go)
* [Python](https://github.com/taibende/pyphilote)

## Authentication

Clients authenticate in Philote using [JSON Web Tokens](https://jwt.io), which consist on a JSON payload detailing the read/write permissions a given connection will have. The payload is hashed with a secret known to Philote so that incoming connections can be verified, this way you can generate tokens in your application backend and use them from the browser client without fear.

Clients in different language will provide methods to generate these tokens, for now, the [Go client](https://github.com/pote/philote-go/blob/master/token.go) should be the reference implementation, although you'll notice that it's an extremely simple one so ports to other languages should be trivial to implement provided with a decent JWT library.

For incoming websockets connections, Philote will look to find the authentication token in the `Authorization` header, but since the native browser JavaScript WebSocket API does not provide a way to manipulate the request headers Philote will also look for the `auth` query parameter in case it fails to authenticate using the header option.

## License

Released under MIT License, check LICENSE file for details.
