/*
   SalioMP
   Copyright (C) 2020 Juby210

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	sio "github.com/gsocket-io/golang-socketio"
	"github.com/gsocket-io/golang-socketio/transport"
)

var (
	config cfg
	wd, _  = os.Getwd()
	mod    = path.Join(wd, "..", "mod")
	last   = pos{-18, -18}
	lastl  = 0
)

func main() {
	cfg := "config.json"
	args := os.Args[1:]
	if len(args) >= 1 {
		cfg = args[0]
	}
	if len(args) >= 2 {
		mod = args[1]
	}
	unmarshal(cfg, &config)

	pathes := getPathes()

	if isNotExist(config.SalioPath) {
		os.Mkdir(config.SalioPath, 0644)
	}
	if isNotExist(pathes["p1"]) {
		os.Create(pathes["p1"])
	}
	ioutil.WriteFile(pathes["mdata"], []byte(`{"level":0,"waiting":false}`), 0644)
	ioutil.WriteFile(pathes["mdata2"], []byte(`{"level":0,"waiting":false}`), 0644)
	ioutil.WriteFile(pathes["p2"], []byte("-18\n-18"), 0644)
	if isNotExist(pathes["json"]) {
		log.Print("Downloading json.lua..")
		res, err := http.Get("https://raw.githubusercontent.com/rxi/json.lua/master/json.lua")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer res.Body.Close()

		file, _ := os.Create(pathes["json"])
		defer file.Close()
		io.Copy(file, res.Body)
		log.Print("Done")
	}

	c, err := sio.Dial(
		sio.GetUrl(config.IP, config.Port, config.Secure),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.On(sio.OnConnection, func(_ *sio.Channel) {
		log.Print("Connected")
	})

	c.On("config", func(_ *sio.Channel, cfg map[string]interface{}) {
		if cfg["syncLevels"].(bool) {
			log.Print("Syncing levels..")
			res, err := http.Get(fmt.Sprintf("http://%s:%d/data.salio", config.IP, config.Port))
			if err != nil {
				log.Fatal(err)
				return
			}
			defer res.Body.Close()

			data, _ := ioutil.ReadAll(res.Body)
			orgdata, _ := ioutil.ReadFile(pathes["data"])
			if string(data) == string(orgdata) {
				log.Print("Done")
				return
			}
			ioutil.WriteFile(pathes["data"], data, 0644)

			log.Print("Downloading levels..")
			res, err = http.Get(fmt.Sprintf("http://%s:%d/levels.zip", config.IP, config.Port))
			if err != nil {
				log.Fatal(err)
				return
			}
			defer res.Body.Close()

			file, _ := os.Create(pathes["lzip"])
			defer file.Close()
			io.Copy(file, res.Body)
			if !isNotExist(pathes["levels"]) {
				os.RemoveAll(pathes["levels"])
			}
			os.Mkdir(pathes["levels"], 0644)
			_, err = unzip(pathes["lzip"], pathes["levels"])
			if err != nil {
				log.Fatal(err)
			}
			log.Print("Done")
			os.Remove(pathes["lzip"])
		}
	})

	c.On("level", func(_ *sio.Channel, data interface{}) {
		log.Print(data)
		marshal(pathes["mdata2"], data)
	})

	c.On("playermove", func(_ *sio.Channel, data pos) {
		ioutil.WriteFile(pathes["p2"], []byte(fmt.Sprintf("%d\n%d", data.X, data.Y)), 0644)
		log.Printf("p2; %d | %d", data.X, data.Y)
	})

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			if event, ok := <-watcher.Events; ok {
				if event.Op&fsnotify.Write != 0 {
					if event.Name == pathes["p1"] {
						content, _ := ioutil.ReadFile(pathes["p1"])
						a := strings.Split(strings.ReplaceAll(string(content), "\r", ""), "\n")
						if len(a) >= 2 {
							x, _ := strconv.Atoi(a[0])
							y, _ := strconv.Atoi(a[1])
							if last.X != x || last.Y != y {
								c.Emit("playermove", pos{x, y})
								last.X = x
								last.Y = y
							}
						}
					} else {
						var data map[string]interface{}
						err = unmarshal(pathes["mdata"], &data)
						if err == nil {
							l := int(data["level"].(float64))
							if l != lastl || data["waiting"].(bool) {
								c.Emit("level", data)
								lastl = l
							}
						}
					}
				}
			}
		}
	}()

	err = watcher.Add(pathes["p1"])
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add(pathes["mdata"])
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan bool)
}

func getPathes() map[string]string {
	config.SalioPath = path.Join(config.SalioPath, "mp")
	smp := config.SalioPath
	return map[string]string{
		"p1":     path.Join(smp, "p1.txt"),
		"p2":     path.Join(smp, "p2.txt"),
		"mdata":  path.Join(smp, "data.json"),
		"mdata2": path.Join(smp, "data2.json"),
		"json":   path.Join(smp, "json.lua"),

		"data":   path.Join(mod, "data.salio"),
		"levels": path.Join(mod, "levels"),
		"lzip":   path.Join(mod, "levels.zip"),
	}
}
