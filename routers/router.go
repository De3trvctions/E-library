package routers

import (
	"e-library/controllers"
	book "e-library/controllers/booking"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Register the routes for the book operations
	web.Router("/Book", &book.BookController{}, "get:Book")
	web.Router("/Borrow", &book.BookController{}, "post:Borrow")
	web.Router("/Extend", &book.BookController{}, "post:Extend")
	web.Router("/All", &book.BookController{}, "get:All")
	web.Router("/Return", &book.BookController{}, "post:Return")

	web.Router("/", &controllers.MainController{})
}
