package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:8000/echo"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net/http"
// )

// func main() {

// 	resp, err := http.Get("http://localhost:8000")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	fmt.Println("Response status:", resp.Status)

// 	scanner := bufio.NewScanner(resp.Body)
// 	for i := 0; scanner.Scan() && i < 5; i++ {
// 		fmt.Println(scanner.Text())
// 	}

// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}
// }
