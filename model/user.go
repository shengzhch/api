package model

type User struct {
	ID     int64
	Name   string
	CanUse int16 `xorm:"default(1)"`
}
