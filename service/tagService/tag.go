package tagService

import (
	"encoding/json"
	"gin_example/models"
	"gin_example/pkg/export"
	"gin_example/pkg/file"
	"gin_example/pkg/gredis"
	"gin_example/pkg/logging"
	"gin_example/service/cache_service"
	"io"
	"path"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"github.com/tealeg/xlsx"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistTagByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

// 先从redis判断是否存在，若不存在再去slow db查
func (t *Tag) ExistTagByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

// 删除对应ID的Tag
func (t *Tag) DeleteTag() (bool, error) {
	return models.DeleteTag(t.ID)
}

// 获取所有符合条件的Tag
func (t *Tag) GetTagTotal(mps map[string]interface{}) (int, error) {
	return models.GetTagTotal(mps)
}

// 新增标签
func (t *Tag) AddTag() bool {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

// 修改标签
func (t *Tag) Edit() (bool, error) {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	if t.Name != "" {
		data["name"] = t.Name
	}
	if t.State != -1 {
		data["state"] = t.State
	}
	return models.EditTag(t.ID, data)
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)
	// 需要和cache_service中的GetTagsKey要使用的字段对应
	cache := cache_service.Tag{
		State: t.State,
		Name:  t.Name,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal([]byte(data), &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.GetPageSize(), t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) Count() (int, error) {
	return models.GetArticleTotal(t.getMaps())
}

// 根据Tag中设定的条件将Tag导出
func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("标签信息")
	if err != nil {
		return "", err
	}
	titles := []string{"ID", "名称", "创建人", "修改人", "修改时间"}
	row := sheet.AddRow()
	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		row = sheet.AddRow()
		for _, val := range values {
			cell = row.AddCell()
			cell.Value = val
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + time + ".xlsx"

	fullPath := path.Join(export.GetExcelFullPath(), filename)
	file.IsNotExistMkDir(export.GetExcelFullPath())
	err = xlsxFile.Save(fullPath)
	if err != nil {
		logging.WarnF("failed to export tags err: %v", err)
		return "", err
	}

	return filename, nil
}

func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := xlsx.GetRows("标签信息")
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}
			models.AddTag(data[1], 1, data[2])
		}
	}

	return nil
}

// 将Tag中的有效字段映射出去
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}

func (t *Tag) GetPageSize() int {
	return max(1, t.PageSize)
}
