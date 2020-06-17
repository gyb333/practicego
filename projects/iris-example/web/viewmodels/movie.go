package viewmodels

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"iris-example/datamodels"
)

type Movie struct {
	datamodels.Movie
}

func (m Movie) IsValid() bool {
	/* do some checks and return true if it's valid... */
	return m.ID > 0
}

func (m Movie) Dispatch(ctx iris.Context) {
	if !m.IsValid() {
		ctx.NotFound()
		return
	}
	ctx.JSON(m, context.JSON{Indent: " "})
}