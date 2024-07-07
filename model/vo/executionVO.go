package vo

import (
	"time"
)

type ExecutionVO struct {
	Id            uint      `json:"id"`            // id
	CreateAt      time.Time `json:"createAt"`      // 代码执行时间戳
	Code          string    `json:"code"`          // 代码
	Language      string    `json:"language"`      // 编程语言
	MaxMemoryCost uint64    `json:"maxMemoryCost"` // 单个输入用例的最大内存消耗
	MaxTimeCost   int64     `json:"maxTimeCost"`   // 单个输入用例的最大时间消耗
	InputList     []string  `json:"inputList"`     //所有的输入用例
	OutputList    []string  `json:"outputList"`    //所有的输入用例
	Status        int8      `json:"status"`        // 执行状态
	KeyPairId     uint      `json:"keyPairId"`     // 使用的密钥对ID
}
