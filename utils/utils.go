package utils

import (
	"github.com/seefs001/seng"
	"github.com/seefs001/seng/middlewares/cors"
	"github.com/seefs001/seng/middlewares/recovery"
)

func ApplyDefaultMiddlewares(e *seng.Engine) {
	e.Use(cors.Default())
	e.Use(recovery.Default())
}
