package controller

import (
	"test/operator-study/cmdbdemo/pkg/controller/cmdbservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, cmdbservice.Add)
}
