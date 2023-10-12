package main

import (
	"fmt"
	"github.com/andyzhou/tinypage"
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
	TplPath = "./tpl"
	StaticPath = "./html"

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
	page := tinypage.NewPage()

	//set core path
	page.SetCorePath(TplPath, StaticPath)

	//register auto gen
	page.RegisterAutoGen("test", 10, nil)

	//add tpl ext func
	page.AddExtFunc("html", nil)

	//add shared tpl
	page.AddSharedTpl("test.tpl")

	//gen page
	page.GenPage("test.tpl", "test", "test.html", nil)

	//start wait group
	wg.Add(1)
	fmt.Println("start example...")

	wg.Wait()
	page.Quit()
	fmt.Println("stop example...")
}