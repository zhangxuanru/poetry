/*
@Time : 2019/9/19 16:46
@Author : zxr
@File : index
@Software: GoLand
*/
package controllers

import (
	"math"
	"net/http"
	"poetry/app/bootstrap"
	"poetry/app/logic"
	"poetry/app/models"
	"poetry/config/define"
	templateHtml "poetry/libary/template"
	"sort"
	"strconv"
)

//首页
func Index(writer http.ResponseWriter, request *http.Request) {
	var (
		err             error
		contentData     define.ContentAll
		categoryData    []models.Category
		famousData      []models.Category
		authorData      []models.Author
		ancientBookData []models.AncientBook
		html            *templateHtml.Html
		assign          map[string]interface{}
		currPage        int    //当前页数
		offset          int    //偏移量
		limit           = 8    //推荐诗词每页显示多少条
		countPage       = 0    //总页数
		pageStr         string //URL传过来的当前页数
	)
	if pageStr = request.FormValue("page"); len(pageStr) == 0 {
		pageStr = "1"
	}
	html = templateHtml.NewHtml(writer)
	currPage, _ = strconv.Atoi(pageStr)
	//推荐总数
	countNum := logic.NewRecommendLogic().GetRecommendCount()
	countPage = int(math.Ceil(float64(countNum) / float64(limit)))
	if currPage == 0 || currPage > countPage || currPage > 100 {
		currPage = 1
	}
	offset = (currPage - 1) * limit
	//获取推荐数据
	if contentData, err = logic.NewRecommendLogic().GetSameDayRecommendPoetryData(offset, limit); err != nil {
		html.DisplayErrorPage(err)
		return
	}
	//获取诗文分类
	if categoryData, err = logic.NewCategoryLogic().GetCateByPositionLimit(define.PoetryShowPosition, 0, 76); err != nil {
		html.DisplayErrorPage(err)
		return
	}
	sort.Slice(categoryData, func(i, j int) bool {
		return len(categoryData[i].CatName) > len(categoryData[j].CatName)
	})
	//获取诗文名句分类
	if famousData, err = logic.NewCategoryLogic().GetCateByPositionLimit(define.FamousShowPosition, 0, 12); err != nil {
		html.DisplayErrorPage(err)
		return
	}
	//获取作者列表
	if authorData, err = logic.NewAuthorLogic().GetListByOrderCountDesc(0, 53); err != nil {
		html.DisplayErrorPage(err)
		return
	}
	//获取古籍目录列表
	if ancientBookData, err = logic.NewAncientBook().GetBookListByLimit(0, 32); err != nil {
		html.DisplayErrorPage(err)
		return
	}
	assign = make(map[string]interface{})
	assign["contentData"] = contentData
	assign["cdnDomain"] = bootstrap.G_Conf.CdnStaticDomain
	assign["webDomain"] = bootstrap.G_Conf.WebDomain
	assign["currPage"] = currPage
	assign["nextPage"] = currPage + 1
	assign["prevPage"] = currPage - 1
	assign["countPage"] = countPage
	assign["categoryData"] = categoryData
	assign["authorData"] = authorData
	assign["famousData"] = famousData
	assign["ancientBookData"] = ancientBookData
	assign["title"] = define.WebTitle
	assign["description"] = define.WebDescription
	assign["urlPath"] = request.URL.Path
	assign["version"] = define.StaticVersion
	html.Display("index.html", assign)
	return
}
