package vo

import "time"

type KeyPairVO struct {
	ID        uint      `json:"id"`
	AccessKey string    `json:"accessKey"`
	SecretKey string    `json:"secretKey"`
	UserId    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
