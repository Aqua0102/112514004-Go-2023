package main

import (
	"bufio"
	"context"

	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
)

type client chan<- string

var (
	entering      = make(chan client)
	leaving       = make(chan client)
	messages      = make(chan rxgo.Item)
	ObservableMsg = rxgo.FromChannel(messages)
)

func broadcaster() {
	clients := make(map[client]bool)
	MessageBroadcast := ObservableMsg.Observe()
	for {
		select {
		case msg := <-MessageBroadcast:

			for cli := range clients {
				cli <- msg.V.(string)
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientWriter(conn *websocket.Conn, ch <-chan string) {
	for msg := range ch {
		conn.WriteMessage(1, []byte(msg))
	}
}

func wshandle(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "你是 " + who + "\n"
	messages <- rxgo.Of(who + " 來到了現場" + "\n")
	entering <- ch

	defer func() {
		log.Println("disconnect !!")
		leaving <- ch
		messages <- rxgo.Of(who + " 離開了" + "\n")
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		messages <- rxgo.Of(who + " 表示: " + string(msg))
	}
}

func InitObservable() {

	swearWords, err := loadWordsFromFile("swear_word.txt")
	if err != nil {
		log.Fatal("Error loading swear words:", err)
	}

	sensitiveNames, err := loadWordsFromFile("sensitive_name.txt")
	if err != nil {
		log.Fatal("Error loading sensitive names:", err)
	}

	ObservableMsg = ObservableMsg.
		Filter(filterSwearWords(swearWords)).
		Map(mapSensitiveNames(sensitiveNames))
}

func mapSensitiveNames(sensitiveNames []string) rxgo.Func {
	return func(_ context.Context, i interface{}) (interface{}, error) {
		message := i.(string)
		for _, name := range sensitiveNames {
			message = strings.Replace(message, name, maskName(name), 1)
		}
		return message, nil
	}
}

func filterSwearWords(swearWords []string) rxgo.Predicate {
	return func(item interface{}) bool {
		message := item.(string)
		for _, word := range swearWords {
			if strings.Contains(message, word) {
				return false
			}
		}
		return true
	}
}

func loadWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func maskName(name string) string {
	runes := []rune(name)
	if len(runes) > 1 {
		if len(runes) > 2 || len(runes) != 2 {
			return string(runes[0]) + strings.Repeat("*", len(runes)-2) + string(runes[len(runes)-1])
		}
		return string(runes[0]) + "*"
	}
	return name
}

func main() {
	InitObservable()
	go broadcaster()
	http.HandleFunc("/wschatroom", wshandle)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("server start at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
