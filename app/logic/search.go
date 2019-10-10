/*
@Time : 2019/10/9 19:33
@Author : zxr
@File : search
@Software: GoLand
*/
package logic

import (
	"poetry/app/models"
	"poetry/config/define"
)

type Searcher interface {
	GetPoetryListByFilter(cstr string, offset, limit int) ([]models.Content, error)
}

type SearchLogic struct {
}

func NewSearchLogic() *SearchLogic {
	return &SearchLogic{}
}

const (
	searchTag     = "tag"     //按tag标签搜索
	searchDynasty = "dynasty" //按朝代搜索
	searchAuthor  = "author"  //按作者搜索
)

//诗文搜索
func (s *SearchLogic) GetSearchShiWenPoetryList(typeStr, cstr string, offset, limit int) (contentData define.ContentAll, err error) {
	var (
		searchMod  Searcher
		poetryList []models.Content
		authorIds  []int64               //作者ID集合
		poetryIds  []int64               //诗词ID集合
		authorData map[int]models.Author //作者信息集合
		tags       TagMp                 //诗词的分类标签信息
	)
	switch typeStr {
	case searchAuthor:
		searchMod = NewAuthorLogic()
	case searchTag:
		searchMod = NewCategoryLogic()
	case searchDynasty:
		searchMod = NewDynastyLogic()
	default:
		return
	}
	if poetryList, err = searchMod.GetPoetryListByFilter(cstr, offset, limit); err != nil || len(poetryList) == 0 {
		return
	}
	poetryIds = ExtractPoetryIdTo64(poetryList)
	authorIds = ExtractAuthorId(poetryList)
	//根据作者ID查询作者表数据
	if authorData, err = NewAuthorLogic().GetAuthorInfoByIds(authorIds); err != nil {
		return
	}
	//根据诗词ID查询分类标签表数据
	tags, _ = NewContentTagLogic().GetDataByPoetryId(poetryIds)
	//将诗词数据，作者数据，朝代数据,分类整合一起
	contentData = NewContentLogic().ProcContentAuthorTagData(poetryList, authorData, tags)
	return
}
