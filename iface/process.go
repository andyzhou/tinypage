package iface

/*
 * interface for lazy process
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

type IProcess interface {
	Quit()
	GenPage(tplFile, subDir, pageFile string,
		dataMap map[string]interface{},
	) bool
	GetTplFace() ITpl
	RegisterAutoGen(tag string, rate int, cb func()) bool
	SetCallBack(cb func(pageFile string, pageData []byte) bool) bool
}