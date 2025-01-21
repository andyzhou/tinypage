package face

import (
	"errors"
	"fmt"
	"github.com/andyzhou/tinypage/iface"
	"log"
	"os"
)

/*
 * face for static page
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//face info
type Static struct {
	staticPath string
}

//construct
func NewStatic(staticPath string) *Static {
	//self init
	this := &Static{
		staticPath:staticPath,
	}
	return this
}

//generate static page
func (f *Static) GenPage(
		tplFile,
		subDir,
		pageFile string,
		dataMap map[string]interface{},
		tplFace iface.ITpl,
	) ([]byte, error) {
	//basic check
	if tplFile == "" || pageFile == "" {
		return nil, errors.New("invalid parameter")
	}
	if dataMap == nil || tplFace == nil {
		return nil, errors.New("invalid parameter")
	}

	//check or create sub dir
	if subDir != "" {
		err := f.checkOrCreateDir(subDir)
		if err != nil {
			return nil, err
		}
	}

	//begin generate page
	pageData, err := tplFace.GenOnePage(tplFile, subDir, pageFile, dataMap)
	if err != nil {
		log.Println("Page::GenStaticPage failed, err:", err.Error())
		return nil, err
	}
	return pageData, nil
}

///////////////
//private func
///////////////

//check or create page dir
func (f *Static) checkOrCreateDir(subDir string) error {
	if subDir == "" {
		return errors.New("invalid sub dir parameter")
	}

	//check or create
	subDirPath := fmt.Sprintf("%s/%s", f.staticPath, subDir)
	err := f.checkOrCreateOneDir(subDirPath)
	if err != nil {
		log.Println("PageFace::checkOrCreateDir failed, err:", err.Error())
		return err
	}
	return nil
}

//check or create dir
func (f *Static) checkOrCreateOneDir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return err
	}
	bRet := os.IsExist(err)
	if bRet {
		return nil
	}
	err = os.Mkdir(dir, 0777)
	return err
}
