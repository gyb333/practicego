package viewmodels

import (

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"iris-example/datamodels"
)

type User struct {
	datamodels.User
}

func (m User) IsValid() bool {
	/* do some checks and return true if it's valid... */
	return m.ID > 0
}


func (m User) Dispatch(ctx iris.Context) {
	if !m.IsValid() {
		ctx.NotFound()
		return
	}
	ctx.JSON(m, context.JSON{Indent: " "})
}
