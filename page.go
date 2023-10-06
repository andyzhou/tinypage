package tinypage

import (
	"errors"
	"github.com/andyzhou/tinypage/face"
	"github.com/andyzhou/tinypage/iface"
)

/*
 * face for page service, main entry
 * this is main entry
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//face info
type Page struct {
	process iface.IProcess
}

//construct
func NewPage() *Page {
	//self init
	this := &Page{
		process:face.NewProcess(),
	}
	return this
}

//quit
func (f *Page) Quit() {
	f.process.Quit()
}

//generate static page
func (f *Page) GenPage(
				tplFile string,
				subDir string,
				pageFile string,
				dataMap map[string]interface{},
			) error {
	return f.process.GenPage(tplFile, subDir, pageFile, dataMap)
}

//register auto gen page func
func (f *Page) RegisterAutoGen(
					tag string,
					rate int,
					cb func(),
				) error {
	return f.process.RegisterAutoGen(tag, rate, cb)
}

//add tpl ext func
func (f *Page) AddExtFunc(
				tag string,
				fun func(arg ...interface{})interface{},
			) bool {
	//get tpl face
	tplFace := f.process.GetTplFace()
	if tplFace == nil {
		return false
	}
	return tplFace.AddExtFunc(tag, fun)
}

//add shared tpl
func (f *Page) AddSharedTpl(
					tplFile string,
				) error {
	//get tpl face
	tplFace := f.process.GetTplFace()
	if tplFace == nil {
		return errors.New("tpl face not init")
	}
	return tplFace.AddSharedTpl(tplFile)
}

//setup core path
func (f *Page) SetCorPath(tplPath, pagePath string) error {
	//check
	if tplPath == "" || pagePath == "" {
		return errors.New("invalid path parameter")
	}
	return f.process.SetCorePath(tplPath, pagePath)
}