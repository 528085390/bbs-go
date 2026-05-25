package main

import (
	"temp/common/db"

	"github.com/zeromicro/go-zero/core/logx"
)

func init() {
	db.Init()
}

func main() {
	logx.Info("common initialized")
}
