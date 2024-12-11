package models

type Article struct {
	Model
	// 嵌套Tag，表示TagID与Tag模型相互关联，在执行查询的时候
	// 能达到Article、Tag关联查询的功能
	TagId int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 通过嵌入的Model对应的回调函数实现了
// func (article *Article) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("CreatedOn", time.Now().Unix())

// 	return nil
// }

// func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
// 	scope.SetColumn("ModifiedOn", time.Now().Unix())

// 	return nil
// }

func CleanAllArticle() bool {
	// Unscoped 表示硬删除
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{})

	return true
}

// 获取单个文章
func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).Find(&article)
	db.Model(&Article{}).Related(&article.Tag)

	return
}

// 获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	// 使用Preload就是一个预加载器，它会执行两条SQL语句
	// 即先查 blog_article; 然后查 blog_tag 然后将查出的结果填充到Article中
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

// 获取文章数
func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)

	return
}

// 根据ID判断文章是否存在
func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)

	return article.ID > 0
}

// 修改文章
func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

// TODO CreatedBy 字段相关的还没有校验

// 新增文章
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagId:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})

	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(&Article{})

	return true
}
