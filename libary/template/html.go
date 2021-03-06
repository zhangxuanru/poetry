/*
@Time : 2019/9/23 19:06
@Author : zxr
@File : html
@Software: GoLand
*/
package template

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"poetry/config/define"
	"strings"
	"unicode/utf8"
)

type Html struct {
	Writer        io.Writer
	templateFiles []string
	ViewPath      string
}

func NewHtml(writer io.Writer) *Html {
	return &Html{
		Writer:   writer,
		ViewPath: define.BaseDir + "/app",
	}
}

//显示错误页面
func (h *Html) DisplayErrorPage(err error) {
	if err == nil {
		err = errors.New("未知错误")
	}
	errStr := fmt.Errorf("%v", err).Error()
	logrus.Println("err:", errStr)
	h.Display404()
	//assign := make(map[string]interface{})
	//assign["errStr"] = errStr
	//h.displaySinglePage("error.html", assign)
	return
}

//显示404页面
func (h *Html) Display404() {
	h.displaySinglePage("404.html", nil)
	return
}

//追加模板
func (h *Html) AddTemplate(template string) {
	h.templateFiles = append(h.templateFiles, template)
}

//解析字符串为html
func unescaped(x string) template.HTML {
	return template.HTML(x)
}

//截取字符串
func subLen(str string, subLen int) string {
	if subLen > utf8.RuneCountInString(str) {
		//subLen = utf8.RuneCountInString(str)
		return str
	}
	content := []rune(str)[0:subLen]
	return string(content)
}

//获取字符串长度
func StrLen(str string) int {
	return utf8.RuneCountInString(str)
}

//去除某些特殊符号
func TrimSymbol(str string) string {
	return strings.Trim(str, "：")
}

//显示模板页
func (h *Html) Display(page string, data map[string]interface{}) {
	htmlPath := []string{
		h.ViewPath + "/views/public/header.html",
		h.ViewPath + "/views/public/naver.html",
		h.ViewPath + "/views/public/footer.html",
		h.ViewPath + "/views/" + page,
	}
	pageName := page
	if strings.Contains(page, "/") {
		index := strings.LastIndex(page, "/")
		pageName = page[index+1:]
	}
	if len(h.templateFiles) > 0 {
		h.templateFiles = append(h.templateFiles, htmlPath...)
	} else {
		h.templateFiles = htmlPath
	}
	tpl := template.New(pageName)
	//切记：加的自定义函数在Parse之前
	tpl = tpl.Funcs(template.FuncMap{"unescaped": unescaped, "subLen": subLen, "strLen": StrLen, "trimSymbol": TrimSymbol})
	tpl = template.Must(tpl.ParseFiles(
		h.templateFiles...,
	))
	if err := tpl.ExecuteTemplate(h.Writer, pageName, data); err != nil {
		logrus.Infoln("ExecuteTemplate error:", err)
	}
	return
}

//显示单页面
func (h *Html) displaySinglePage(page string, data interface{}) {
	file := h.ViewPath + "/views/" + page
	must := template.Must(template.New(page).Funcs(template.FuncMap{"unescaped": unescaped}).ParseFiles(file))
	if err := must.ExecuteTemplate(h.Writer, page, data); err != nil {
		logrus.Infoln("ExecuteTemplate error:", err)
	}
	return
}
