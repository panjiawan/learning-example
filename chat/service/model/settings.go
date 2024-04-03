package model

// Setting 系统设置
type Setting struct {
	ID      int    `gorm:"autoIncrement:true;primaryKey;column:id;type:int;not null"`
	Base    string `gorm:"column:base;type:mediumtext;default:null;comment:'基本设置'"`    // 基本设置
	Live    string `gorm:"column:live;type:mediumtext;default:null;comment:'直播设置'"`    // 直播设置
	Finance string `gorm:"column:finance;type:mediumtext;default:null;comment:'财务设置'"` // 财务设置
	Badge   string `gorm:"column:badge;type:mediumtext;default:null;comment:'图标设置'"`   // 图标设置
	Other   string `gorm:"column:other;type:mediumtext;default:null;comment:'其他设置'"`   // 其他设置
	Version string `gorm:"column:version;type:mediumtext;default:null;comment:'版本设置'"` // 版本设置
}

// SettingColumns get sql column name.获取数据库列名
var SettingColumns = struct {
	ID      string
	Base    string
	Live    string
	Finance string
	Badge   string
	Other   string
	Version string
}{
	ID:      "id",
	Base:    "base",
	Live:    "live",
	Finance: "finance",
	Badge:   "badge",
	Other:   "other",
	Version: "version",
}

type SettingModel struct {
}

func (Setting) TableName() string {
	return "settings"
}

func (a *SettingModel) Get() (*Setting, error) {
	res := &Setting{}
	db := mysqlHandles["live"].Handle().Where("id>0").First(res)
	if db.Error != nil {
		return nil, db.Error
	}

	return res, nil
}
