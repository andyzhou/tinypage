package iface

/*
 * interface for tpl engine
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

type ITpl interface {
	GenOnePage(mainTplFile, subDir, pageFile string, data interface{}) ([]byte, error)
	ResetSharedTpl()
	AddSharedTpl(tplFile string) error
	AddExtFunc(tag string, fun interface{}) bool
}