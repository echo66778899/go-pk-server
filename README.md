# Poker Game Server
Welcome to the Poker Game Server project! This server is designed to handle the core functionality of a poker game, including player connections, game management, and hand evaluation. The server is implemented in Golang and provides a robust and scalable solution for hosting poker games.

## Features
- Golang-Powered: The server is built using Golang, known for its performance and concurrency capabilities, making it ideal for handling multiple simultaneous connections.
- Event-Driven Architecture: The server supports an event-driven model, allowing for efficient handling of game events and player actions.
- Brute Force Hand Evaluation: Implements a brute force method for evaluating poker hands, ensuring accurate and reliable results.
- WebSocket Communication: Supports WebSocket connections for real-time communication between the server and clients. The server will be upgraded to support secure WebSocket (WSS) connections in the near future.
- JSON Data Format: All data is presented and synchronized using JSON, ensuring compatibility and ease of use across different platforms and languages.

## Getting Started
### Prerequisites
Go 1.23 or higher: Ensure you have Go installed. You can download it from `golang.org`.
### Installation
Clone the repository to your local machine:

```
git clone https://github.com/haiphan1811/go-pk-server.git
cd go-pk-server
```

Install the required dependencies:
```
go mod tidy
```
Configuration
Before running the server, you may need to configure certain parameters such as server port, WebSocket settings, and other options. These can be adjusted in the config.json file or through environment variables.

Generate code
```
protoc --go_out=. message.proto
```

### Usage
To start the Poker Game Server:

```
go run main.go
```
This will start the server and listen for incoming WebSocket connections on the specified port.

### WebSocket API
Clients can connect to the server using WebSocket and communicate using JSON-formatted messages. Here's an example of a typical message exchange:

Client Request:

```
{
    "action": "join_game",
    "player_id": "player123",
    "game_id": "game456"
}
```
Server Response:

```
{
    "event": "game_joined",
    "player_id": "player123",
    "game_id": "game456",
    "status": "success"
}
```
Example
Here’s a simple example of how to connect to the server and join a game:
```
package main

import (
    "log"
    "golang.org/x/net/websocket"
)

func main() {
    origin := "http://localhost/"
    url := "ws://localhost:8080/ws"
    ws, err := websocket.Dial(url, "", origin)
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    msg := `{"action": "join_game", "player_id": "player123", "game_id": "game456"}`
    if _, err := ws.Write([]byte(msg)); err != nil {
        log.Fatal(err)
    }

    var response = make([]byte, 512)
    if n, err := ws.Read(response); err == nil {
        log.Printf("Received: %s", response[:n])
    }
}
```

### Testing
To run the tests, use the following command:

```
go test ./...
```
This will execute all the test cases to ensure the server operates as expected.

### Documentation

### License
This project is licensed under the GPL 3.0 License - see the LICENSE file for details.

### Contact
For any questions, feel free to reach out to me at [haiphan1811.dev@google.com].
