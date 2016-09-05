// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

var url map[string]func() = map[string]func(){
	"booktxt": booktxt,
	"45xs":    xs45,
}

var dateFile = "/webser/logs/ztj.date.b"

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("watch-modified file:", event.Name)
					for _, f := range url {
						f()
					}
				}

			case err := <-watcher.Errors:
				log.Println("watch-error:", err)
			}
		}
	}()

	err = watcher.Add(dateFile)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
