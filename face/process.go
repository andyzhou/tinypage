package face

import (
	"github.com/andyzhou/tinyPage/define"
	"github.com/andyzhou/tinyPage/iface"
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
	reqChan chan tinyPageReq
	closeChan chan bool
	autoChan chan bool
}

//construct
func NewProcess(
			tplPath string,
			staticPath string,
		) *Process {
	//self init
	this := &Process{
		tpl:NewTpl(tplPath, staticPath),
		static:NewStatic(staticPath),
		autoMap:make(map[string]tinyAutoGen),
		tickerMap:new(sync.Map),
		reqChan:make(chan tinyPageReq, define.TinyPageChanSize),
		closeChan:make(chan bool, 1),
		autoChan:make(chan bool, 1),
	}
	//spawn main process
	go this.runMainProcess()

	//spawn auto gen process
	go this.runAutoGenProcess()
	return this
}

//quit
func (f *Process) Quit() {
	f.closeChan <- true
	f.autoChan <- true
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
				) (bRet bool) {
	//basic check
	if tplFile == "" || pageFile == "" || dataMap == nil {
		bRet = false
		return
	}

	//try catch panic
	defer func() {
		if err := recover(); err != nil {
			log.Println("Process::GenPage panic, err:", err)
			bRet = false
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
	bRet = true
	return
}

//register auto generate page func
func (f *Process) RegisterAutoGen(
					tag string,
					rate int,
					cb func(),
				) bool {
	//basic check
	if tag == "" || cb == nil {
		return false
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
	return true
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


///////////////
//private func
///////////////

//run main process
func (f *Process) runMainProcess() {
	var (
		req tinyPageReq
		needQuit, isOk bool
	)

	//loop
	for {
		if needQuit && len(f.reqChan) <= 0 {
			break
		}
		select {
		case req, isOk = <- f.reqChan:
			if isOk {
				f.genPageProcess(&req)
			}
		case <- f.closeChan:
			needQuit = true
		}
	}

	//close chan
	close(f.reqChan)
	close(f.closeChan)
}

//run auto gen process
func (f *Process) runAutoGenProcess() {
	var (
		ticker = time.NewTicker(time.Second * define.TinyPageAutoGenRate)
		needQuit bool
	)
	//loop
	for {
		if needQuit {
			break
		}
		select {
		case <- ticker.C:
			f.autoGenProcess()
		case <- f.autoChan:
			needQuit = true
		}
	}
	//close chan
	close(f.autoChan)
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
				) bool {
	//basic check
	if req == nil {
		return false
	}

	//try generate page file
	pageData, bRet := f.static.GenPage(
							req.tplFile,
							req.subDir,
							req.pageFile,
							req.dataMap,
							f.tpl,
						)
	if !bRet {
		return false
	}

	//run call back
	if f.cb != nil {
		f.cb(req.pageFile, pageData)
	}

	return true
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