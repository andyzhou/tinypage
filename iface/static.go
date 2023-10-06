package iface

/*
 * interface for static page
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

type IStatic interface {
	GenPage(
		tplFile,
		subDir,
		pageFile string,
		dataMap map[string]interface{},
		tplFace ITpl,
	) ([]byte, error)
}