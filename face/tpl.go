package face

import (
	"fmt"
	"github.com/andyzhou/tinyPage/define"
	"github.com/pkg/errors"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
 * face for tpl engine
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//face info
type Tpl struct {
	tplPath string
	staticPath string
	extFuncMap map[string]interface{}
	subTpl []string
}

//construct
func NewTpl(
			tplPath string,
			staticPath string,
		) *Tpl {
	//self init
	this := &Tpl{
		tplPath:tplPath,
		staticPath:staticPath,
		extFuncMap:make(map[string]interface{}),
		subTpl:make([]string, 0),
	}
	return this
}

//generate one static page
func (f *Tpl) GenOnePage(
				mainTplFile string,
				subDir string,
				pageFile string,
				data interface{},
			) ([]byte, error) {
	var (
		err error
	)

	//basic check
	if mainTplFile == "" || pageFile == "" || data == nil {
		return nil, errors.New("invalid parameter")
	}

	//parse tpl files
	tpl, err := f.parse(mainTplFile)
	if err != nil {
		return nil, err
	}

	//generate static page
	return f.genStaticPage(subDir, pageFile, tpl, data)
}

//add sub tpl
func (f *Tpl) AddSubTpl(
				tplFile string,
			) bool {
	//basic check
	if tplFile == "" {
		return false
	}
	//add into running sub tpl slice
	tplFile = fmt.Sprintf("%s/%s", f.tplPath, tplFile)
	f.subTpl = append(f.subTpl, tplFile)
	return true
}

//add extend func
//used for dynamic tpl func
func (f *Tpl) AddExtFunc(
				tag string,
				fun interface{},
			) bool {
	//basic check
	if tag == "" || fun == nil {
		return false
	}
	_, ok := f.extFuncMap[tag]
	if ok {
		return false
	}

	//add into running ext map
	f.extFuncMap[tag] = fun
	return true
}

///////////////
//private func
///////////////

//generate static page file
func (f *Tpl) genStaticPage(
				subDir string,
				pageFile string,
				tpl *template.Template,
				data interface{},
			) ([]byte, error) {
	var (
		pageFilePath string
	)

	//format page file path
	if subDir != "" {
		pageFilePath = fmt.Sprintf("%s/%s/%s%s",
							f.staticPath,
							subDir,
							pageFile,
						)
	}else{
		pageFilePath = fmt.Sprintf("%s/%s%s",
							f.staticPath,
							pageFile,
						)
	}

	//create page
	out, err := os.Create(pageFilePath)
	if err != nil {
		return nil, err
	}

	//output page file
	defer out.Close()
	err = tpl.Execute(out, data)
	if err != nil {
		return nil, err
	}

	//read page file
	byteData, err := ioutil.ReadFile(pageFilePath)

	return byteData, err
}

//parse tpl
func (f *Tpl) parse(
				mainTpl string,
			) (*template.Template, error) {
	//init template
	tpl := template.New(mainTpl)

	//add extend function
	f.addFuncMap(tpl)

	//format relate tpl files
	mainTpl = fmt.Sprintf("%s/%s", f.tplPath, mainTpl)

	//common tpl files
	commonTplFiles := make([]string, 0)
	commonTplFiles = append(commonTplFiles, f.subTpl...)
	commonTplFiles = append(commonTplFiles, mainTpl)

	//parse tpl file
	tpl, err := tpl.ParseFiles(commonTplFiles...)
	if err != nil {
		log.Println("Tpl::parse failed, error!" + err.Error())
		return nil ,err
	}
	return tpl, nil
}


//add extend function map
func (f *Tpl) addFuncMap(
				tpl *template.Template,
			) bool {
	if f.extFuncMap == nil || len(f.extFuncMap) <= 0 {
		return false
	}
	funcMap := template.FuncMap{}
	for k, v := range f.extFuncMap {
		funcMap[k] = v
	}
	tpl = tpl.Funcs(funcMap)
	return true
}

//convert time stamp to date
func (f *Tpl) timeStamp2Date(timeStamp int64) string {
	dateTime := time.Unix(timeStamp, 0).Format(define.TimeLayOut)
	tempSlice := strings.Split(dateTime, " ")
	if tempSlice == nil || len(tempSlice) <= 0 {
		return ""
	}
	return tempSlice[0]
}

//convert timestamp like 'Oct 10, 2020' format
func (f *Tpl) timeStampToDayStr(timeStamp int64) string {
	date := f.timeStamp2Date(timeStamp)
	if date == "" {
		return  ""
	}
	tempSlice := strings.Split(date, "-")
	if tempSlice == nil || len(tempSlice) < 3 {
		return ""
	}
	year := tempSlice[0]
	month, _ := strconv.Atoi(tempSlice[1])
	day := tempSlice[2]
	return fmt.Sprintf("%s %s, %s", time.Month(month).String(), day, year)
}

//convert timestamp to data time string format
func (f *Tpl) timeStamp2DateTime(timeStamp int64) string {
	return time.Unix(timeStamp, 0).Format(define.TimeLayOut)
}

///////////////////
//inter tpl ext func
///////////////////

//add inter tpl ext func
func (f *Tpl) addInterExtFunc() {
	f.AddExtFunc("dateTime", f.funcOfDateTime)
	f.AddExtFunc("dayTime", f.funcOfDayTime)
	f.AddExtFunc("date", f.funcOfDate)
	f.AddExtFunc("html", f.funcOfHtml)
}

func (f *Tpl) funcOfDateTime(timeStamp int64) string {
	var (
		dateTime string
	)
	if timeStamp <= 0 {
		return dateTime
	}
	return f.timeStamp2DateTime(timeStamp)
}

//like 'Oct 10, 2020' format
func (f *Tpl) funcOfDayTime(timeStamp int64) string {
	var (
		dateTime string
	)
	if timeStamp <= 0 {
		return dateTime
	}
	return f.timeStampToDayStr(timeStamp)
}

//extend function of date convert
//like YYYY-MM-DD
func (f *Tpl) funcOfDate(timeStamp int64) string {
	var (
		dateTime string
	)
	if timeStamp <= 0 {
		return dateTime
	}
	return f.timeStamp2Date(timeStamp)
}

//extend function of html
func (f *Tpl) funcOfHtml(text string) template.HTML {
	return template.HTML(text)
}