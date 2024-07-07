package vo

type PageDataVO struct {
	Total     int64       `json:"total"`     // 总数据数目
	PageCount int64       `json:"pageCount"` // 页码数量
	Data      interface{} `json:"data"`      // 实际数据
}
