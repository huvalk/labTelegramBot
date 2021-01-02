package main

import (
	"github.com/kataras/golog"
	"labTelegramBot/core/server"
	"sync"
)

func main() {
	app, err := server.NewApp()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = app.StartApp()
		if err != nil {
			golog.Fatal("App down: ", err)
		}

		wg.Done()
	}()

	wg.Wait()
}
