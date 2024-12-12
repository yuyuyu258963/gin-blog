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

// 判断该Article是否存在
func (a *Article) ExistArticleByID() (bool, error) {
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		return true, nil
	}

	return models.ExistArticleByID(a.ID)
}
