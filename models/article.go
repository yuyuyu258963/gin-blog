package models

type Article struct {
	Model
	// 嵌套Tag，表示TagID与Tag模型相互关联，在执行查询的时候
	// 能达到Article、Tag关联查询的功能
	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func CleanAllArticle() (bool, error) {
	// Unscoped 表示硬删除
	err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{}).Error

	return err == nil, nil
}

// 获取单个文章
func GetArticle(id int) (article *Article, err error) {
	article = &Article{}
	err = db.Where("id = ? AND deleted_on = ?", id, 0).First(article).Related(&article.Tag).Error

	return
}

// 获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article, err error) {
	// 使用Preload就是一个预加载器，它会执行两条SQL语句
	// 即先查 blog_article; 然后查 blog_tag 然后将查出的结果填充到Article中
	err = db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	// err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error

	return
}

// 获取文章数
func GetArticleTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Article{}).Where(maps).Count(&count).Error
	return
}

// 根据ID判断文章是否存在
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error

	return article.ID > 0 && err == nil, err
}

// 修改文章
func EditArticle(id int, data interface{}) (bool, error) {
	err := db.Model(&Article{}).Where("id = ?", id).Updates(data).Error

	return err == nil, err
}

// TODO CreatedBy 字段相关的还没有校验

// 新增文章
func AddArticle(data map[string]interface{}) (bool, error) {
	err := db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	}).Error

	return err == nil, nil
}

func DeleteArticle(id int) (bool, error) {
	err := db.Where("id = ?", id).Delete(&Article{}).Error

	return err == nil, nil
}
