package articleService

import (
	"encoding/json"
	"gin_example/models"
	"gin_example/pkg/gredis"
	"gin_example/pkg/logging"
	"gin_example/service/cache_service"
	"log"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

// 优先从Redis中取出Article若不存在再去Slow DB取
func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal([]byte(data), &cacheArticle)
			log.Printf("load article from redis: %s", key)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600)
	return article, nil
}

// 删除Article
func (a *Article) Delete() (bool, error) {
	return models.DeleteArticle(a.ID)
}

func (a *Article) Add() (bool, error) {
	maps := a.GetMaps()
	return models.AddArticle(maps)
}

// 编辑文章内容
func (a *Article) Edit() (bool, error) {
	maps := a.GetMaps()
	return models.EditArticle(a.ID, maps)
}

// 获得该记录的条数
func (a *Article) Count() (int, error) {
	maps := a.GetMaps()
	return models.GetArticleTotal(maps)
}

// 获取所有指定范围内的Article
func (a *Article) GetAll() (cacheArticles []models.Article, err error) {
	// 要和cache_service中的GetArticlesKey要使用的字段对应
	cache := cache_service.Article{
		ID:       a.ID,
		TagID:    a.TagID,
		State:    a.State,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()
	cacheArticles = make([]models.Article, 0)
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Warn(err)
		} else {
			json.Unmarshal([]byte(data), &cacheArticles)
			log.Println("log from cache")
			return cacheArticles, nil
		}
	}

	maps := a.GetMaps()
	cacheArticles, err = models.GetArticles(
		a.PageNum,
		a.GetPageSize(),
		maps,
	)

	// fmt.Printf("maps get all %+v err : %v ", cacheArticles, err)
	if err != nil {
		logging.WarnF("GetArticles failed err : %v", err)
		return nil, err
	}

	gredis.Set(key, cacheArticles, 3600)
	return
}

// 判断该Article是否存在
func (a *Article) ExistArticleByID() (bool, error) {
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		return true, nil
	}

	return models.ExistArticleByID(a.ID)
}

// 将其中修改过的字段额外保存到map中
func (a *Article) GetMaps() map[string]interface{} {
	data := make(map[string]interface{}, 0)
	if a.TagID > 0 {
		data["tag_id"] = a.TagID
	}
	if a.Title != "" {
		data["title"] = a.Title
	}
	if a.Desc != "" {
		data["desc"] = a.Desc
	}
	if a.Content != "" {
		data["content"] = a.Content
	}
	if a.CoverImageUrl != "" {
		data["cover_image_url"] = a.CoverImageUrl
	}
	if a.ModifiedBy != "" {
		data["modified_by"] = a.ModifiedBy
	}

	return data
}

func (a *Article) GetPageSize() int {
	return max(1, a.PageSize)
}
