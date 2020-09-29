package main

import (
	"fmt"
	"github.com/andyzhou/tinyPage"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

/*
 * face for example client
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

const (
	TplPath = "/data/tpl"
	StaticPath = "/data/html"

)

func main() {
	var (
		wg sync.WaitGroup
	)

	//try catch signal
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		os.Kill,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)

	///signal snatch
	go func(wg *sync.WaitGroup) {
		var needQuit bool
		for {
			if needQuit {
				break
			}
			select {
			case s := <- c:
				log.Println("Get signal of ", s.String())
				wg.Done()
				needQuit = true
			}
		}
	}(&wg)

	//init api face
	page := tinyPage.NewPage(TplPath, StaticPath)

	//add tpl ext func
	page.AddExtFunc("html", nil)


	//start wait group
	wg.Add(1)
	fmt.Println("start example...")


	wg.Wait()
	page.Quit()
	fmt.Println("stop example...")
}