package model

import (
	"time"
)

type Task struct {
	Id int "json:id"
	Title string "json:title"
	Done bool "json:done"
	CreatedDate time.Time "json:createdDate"
}
