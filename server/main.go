package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	sio "github.com/gsocket-io/golang-socketio"
	"github.com/gsocket-io/golang-socketio/transport"
)

var (
	config map[string]interface{}
	p1, p2 *sio.Channel
)

type pos struct {
	X, Y int
}

func main() {
	parseConfig()

	io := sio.NewServer(transport.GetDefaultWebsocketTransport())

	io.On(sio.OnConnection, func(ch *sio.Channel) {
		if p1 == nil {
			p1 = ch
		} else {
			p2 = ch
		}

		ch.Emit("config", config)
	})

	io.On("playermove", func(ch *sio.Channel, data pos) {
		isp1, prefix := isP1(ch)
		log.Printf("%s; %d | %d", prefix, data.X, data.Y)

		if isp1 {
			if p2 != nil {
				p2.Emit("playermove", data)
			}
		} else {
			if p1 != nil {
				p1.Emit("playermove", data)
			}
		}
	})

	io.On(sio.OnDisconnection, func(ch *sio.Channel) {
		isp1, _ := isP1(ch)
		if isp1 {
			p1 = nil
			if p2 != nil {
				p2.Emit("playermove", pos{-18, -18})
			}
		} else {
			p2 = nil
			if p1 != nil {
				p1.Emit("playermove", pos{-18, -18})
			}
		}
	})

	http.Handle("/socket.io/", io)
	http.HandleFunc("/data.salio", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "data.salio")
	})
	http.HandleFunc("/levels.zip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		http.ServeFile(w, r, "levels.zip")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/juby210-PL/SalioMP", http.StatusMovedPermanently)
	})
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%v", config["port"].(float64)), nil))
}

func parseConfig() {
	byteValue, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal(err)
	}
}

func isP1(ch *sio.Channel) (is bool, prefix string) {
	is = false
	prefix = "p2"
	if ch == p1 {
		is = true
		prefix = "p1"
	}
	return
}
