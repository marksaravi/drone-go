package main

import (
	"github.com/MarkSaravi/drone-go/modules/esc"
	"github.com/MarkSaravi/drone-go/types"
)

func createESCsHandler() types.ESCsHandler {
	return esc.NewESCsHandler()
}
