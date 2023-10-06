package face

import (
	"errors"
	"github.com/andyzhou/tinypage/define"
	"github.com/andyzhou/tinypage/iface"
	"log"
	"sync"
	"time"
)

/*
 * face for lazy process
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//inter request
type tinyPageReq struct {
	tplFile string
	subDir string
	pageFile string
	dataMap map[string]interface{}
}

//inter auto gen
type tinyAutoGen struct {
	rate int
	fun func()
}

//face info
type Process struct {
	tpl iface.ITpl
	static iface.IStatic
	cb func(pageFile string, pageData []byte) bool
	autoMap map[string]tinyAutoGen
	tickerMap *sync.Map
	initDone bool
	reqChan chan tinyPageReq
	closeChan chan bool
	autoCloseChan chan bool
}

//construct
func NewProcess() *Process {
	//self init
	this := &Process{
		autoMap:make(map[string]tinyAutoGen),
		tickerMap:new(sync.Map),
		reqChan:make(chan tinyPageReq, define.TinyPageChanSize),
		closeChan:make(chan bool, 1),
		autoCloseChan:make(chan bool, 1),
	}
	//spawn main process
	go this.runMainProcess()

	//spawn auto gen process
	go this.runAutoGenProcess()
	return this
}

//quit
func (f *Process) Quit() {
	if f.closeChan != nil {
		f.closeChan <- true
	}
	if f.autoCloseChan != nil {
		f.autoCloseChan <- true
	}
}

//get tpl face
func (f *Process) GetTplFace() iface.ITpl {
	return f.tpl
}

//generate static page
func (f *Process) GenPage(
					tplFile string,
					subDir string,
					pageFile string,
					dataMap map[string]interface{},
				) error {
	//basic check
	if tplFile == "" || pageFile == "" || dataMap == nil {
		return errors.New("invalid parameter")
	}
	if !f.initDone {
		return errors.New("core tpl and page path not setup")
	}

	//try catch panic
	defer func() {
		if err := recover(); err != nil {
			log.Println("tinyPage.Process::GenPage panic, err:", err)
			return
		}
	}()

	//init request
	req := tinyPageReq{
		tplFile:tplFile,
		subDir:subDir,
		pageFile:pageFile,
		dataMap:dataMap,
	}

	//send to chan
	f.reqChan <- req
	return nil
}

//register auto generate page func
func (f *Process) RegisterAutoGen(
					tag string,
					rate int,
					cb func(),
				) error {
	//basic check
	if tag == "" || cb == nil {
		return errors.New("invalid parameter")
	}

	//check and set default value
	if rate < define.TinyPageAutoGenRate {
		rate = define.TinyPageAutoGenRate
	}

	//sync into running map
	autoGen := tinyAutoGen{
		rate:rate,
		fun:cb,
	}
	f.autoMap[tag] = autoGen
	f.tickerMap.Store(tag, time.Now().Unix())
	return nil
}

//set callback for gen page success
func (f *Process) SetCallBack(
					cb func(pageFile string, pageData []byte) bool,
				) bool {
	//basic check
	if cb == nil {
		return false
	}
	if f.cb != nil {
		return false
	}
	f.cb = cb
	return true
}

//set core path
func (f *Process) SetCorePath(tplPath, staticPath string) error {
	//check
	if tplPath == "" || staticPath == "" {
		return errors.New("invalid path parameter")
	}
	if f.initDone {
		return errors.New("path had init")
	}
	//init tpl and static obj
	f.tpl = NewTpl(tplPath, staticPath)
	f.static = NewStatic(staticPath)
	f.initDone = true
	return nil
}

///////////////
//private func
///////////////

//run main process
func (f *Process) runMainProcess() {
	var (
		req tinyPageReq
		isOk bool
	)

	//defer
	defer func() {
		if err := recover(); err != nil {
			log.Printf("tinypage.process panic, err:%v\n", err)
		}
		//close chan
		close(f.reqChan)
		close(f.closeChan)
	}()

	//loop
	for {
		select {
		case req, isOk = <- f.reqChan:
			if isOk && &req != nil {
				f.genPageProcess(&req)
			}
		case <- f.closeChan:
			return
		}
	}
}

//run auto gen process
func (f *Process) runAutoGenProcess() {
	var (
		ticker = time.NewTicker(time.Second * define.TinyPageAutoGenRate)
	)

	//defer
	defer func() {
		if err := recover(); err != nil {
			log.Printf("tinyPage.process panic, err:%v\n", err)
		}
		ticker.Stop()
		//close chan
		close(f.autoCloseChan)
	}()

	//loop
	for {
		select {
		case <- ticker.C:
			f.autoGenProcess()
		case <- f.autoCloseChan:
			return
		}
	}
}

//process auto gen page opt
func (f *Process) autoGenProcess() bool {
	var (
		lastTime int64
		diff int64
	)

	//basic check
	if f.autoMap == nil || len(f.autoMap) <= 0 {
		return false
	}

	//get current time
	now := time.Now().Unix()

	//process one by one
	for tag, autoGen := range f.autoMap {
		//get last ticker time
		lastTime = f.getLastTickerTime(tag)
		if lastTime <= 0 {
			//get failed
			continue
		}
		diff = now - lastTime
		if diff < int64(autoGen.rate) {
			//in rate, do nothing
			continue
		}
		//run call cb
		autoGen.fun()
		//sync last ticker
		f.tickerMap.Store(tag, now)
	}

	return true
}

//process page generate opt
func (f *Process) genPageProcess(
					req *tinyPageReq,
				) error {
	//basic check
	if req == nil {
		return errors.New("invalid parameter")
	}

	//generate page file
	pageData, err := f.static.GenPage(
							req.tplFile,
							req.subDir,
							req.pageFile,
							req.dataMap,
							f.tpl,
						)
	if err != nil {
		return err
	}

	//run call back
	if f.cb != nil {
		f.cb(req.pageFile, pageData)
	}

	return nil
}

//get last ticker time
func (f *Process) getLastTickerTime(
					tag string,
				) int64 {
	var (
		lastTime int64
	)

	//basic check
	if tag == "" || f.tickerMap == nil {
		return lastTime
	}
	v, ok := f.tickerMap.Load(tag)
	if !ok {
		return lastTime
	}
	lastTime, ok = v.(int64)
	if !ok {
		return lastTime
	}
	return lastTime
}