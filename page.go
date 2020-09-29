package tinyPage

import (
	"github.com/andyzhou/tinyPage/face"
	"github.com/andyzhou/tinyPage/iface"
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
func NewPage(
		tplPath string,
		staticPath string,
	) *Page {
	//self init
	this := &Page{
		process:face.NewProcess(tplPath, staticPath),
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
			) bool {
	return f.process.GenPage(tplFile, subDir, pageFile, dataMap)
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

//add sub tpl
func (f *Page) AddSubTpl(
					tplFile string,
				) bool {
	//get tpl face
	tplFace := f.process.GetTplFace()
	if tplFace == nil {
		return false
	}
	return tplFace.AddSubTpl(tplFile)
}