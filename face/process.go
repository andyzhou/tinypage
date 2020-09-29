package face

import (
	"github.com/andyzhou/tinyPage/define"
	"github.com/andyzhou/tinyPage/iface"
	"log"
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

//face info
type Process struct {
	tpl iface.ITpl
	static iface.IStatic
	cb func(pageFile string, pageData []byte) bool
	reqChan chan tinyPageReq
	closeChan chan bool
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
		reqChan:make(chan tinyPageReq, define.TinyPageChanSize),
		closeChan:make(chan bool, 1),
	}
	//spawn main process
	go this.runMainProcess()
	return this
}

//quit
func (f *Process) Quit() {
	f.closeChan <- true
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