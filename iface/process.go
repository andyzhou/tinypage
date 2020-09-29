package iface

/*
 * interface for lazy process
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

type IProcess interface {
	Quit()
	GenPage(
		tplFile string,
		subDir string,
		pageFile string,
		dataMap map[string]interface{},
	) bool
	GetTplFace() ITpl
	SetCallBack(cb func(pageFile string, pageData []byte) bool) bool
}