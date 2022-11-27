package iface

/*
 * interface for tpl engine
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

type ITpl interface {
	GenOnePage(mainTplFile, subDir, pageFile string, data interface{}) ([]byte, error)
	AddSubTpl(tplFile string) bool
	AddExtFunc(tag string, fun interface{}) bool
}