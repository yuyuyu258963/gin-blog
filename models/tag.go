package models

import "gin_example/pkg/logging"

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 通过嵌入的Model对应的回调函数实现了

/*
	这两个属于gorm的Callbacks，可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时
	将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改
*/

// 新增两个Hook函数
// func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("CreatedOn", time.Now().Unix())

// 	return nil
// }

// func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
// 	scope.SetColumn("ModifiedOn", time.Now().Unix())

// 	return nil
// }

func CleanAllTag() bool {
	// Unscoped 表示硬删除
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{})

	return true
}

// 获取所有标签列表
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag, err error) {
	err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	return
}

// 获取标签数
func GetTagTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Tag{}).Where(maps).Count(&count).Error
	return
}

// 通过tag名称检查该tag是否存在
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil {
		return false, err
	}
	return tag.ID > 0, nil
}

// 通过TagID检测该tag是否存在
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil {
		logging.InfoF("failed run ExistTagByID %d", id)
		return false, err
	}
	return tag.ID > 0, nil
}

// 新增Tag
func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})

	return true
}

// 修改标签
func EditTag(id int, data map[string]interface{}) (bool, error) {
	err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// 删除一个标签
func DeleteTag(id int) (bool, error) {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
