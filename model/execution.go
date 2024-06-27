package model

import "gorm.io/gorm"

type Execution struct {
	gorm.Model
	Code          string `gorm:"type:longtext;null;comment:代码" json:"code"`
	Language      string `gorm:"type:varchar(20);null;comment:语言" json:"language"`
	MaxMemoryCost uint64 `gorm:"type:bigint;null;comment:单个输入用例的最大内存消耗" json:"maxMemoryCost"`
	MaxTimeCost   int64  `gorm:"type:bigint;null;comment:单个输入用例的最大时间消耗" json:"maxTimeCost"`
	// 注意，该值需要转为 []sting 后才能使用，OutputList 同理
	InputList  string `gorm:"type:longtext;null;comment:输入用例" json:"inputList"`
	OutputList string `gorm:"type:longtext;null;comment:输出用例" json:"outputList"`
	Status     int8   `gorm:"type:tinyint;null;comment:执行状态" json:"status"`
	User       User
	// 虽然更好的方式是维护 keyPairId 即可，但是为了方便，把使用的 userId 也维护了
	UserId    uint `json:"userId"`
	KeyPairId uint `json:"keyPairId"`
}

func (e Execution) TableName() string {
	return "executions"
}
