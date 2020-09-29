package iface

/*
 * interface for static page
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

type IStatic interface {
	GenPage(
		tplFile string,
		subDir string,
		pageFile string,
		dataMap map[string]interface{},
		tplFace ITpl,
	) ([]byte, bool)
}