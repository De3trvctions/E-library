package booking

import (
	"e-library/consts"
	"e-library/controllers/system"
	"e-library/models"
	"e-library/models/dto"
	"e-library/utility"
	"e-library/validation"
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type BookController struct {
	system.PermissionController
}

func (ctl *BookController) Prepare() {
	ctl.PermissionController.Prepare()
}

// All
//
//	@Title			Get all book
//	@Description	Get all book with details
//	@Success		200			{object}	web.M
//	@router			/All [get]
func (ctl *BookController) All() {
	book := models.RentingBook{}
	result, errCode, err := book.GetAll()
	if err != nil || errCode != 0 {
		if errCode == 0 {
			errCode = consts.OPERATION_FAILED
		}
		logs.Error("[BookController][All]Db error:", err)
		ctl.Error(errCode)
	}

	ctl.Success(web.M{"Books": result})
}

// Book
//
//	@Title			Get all book
//	@Description	Get all book with details
//	@Success		200			{object}	web.M
//	@router			/Book [get]
func (ctl *BookController) Book() {
	req := dto.ReqGetBook{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[BookController][Book] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[BookController][Book] FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	book := models.RentingBook{}
	result, err := book.GetBook(req.BookTitle)
	if err != nil {
		logs.Error("[BookController][Book]Db error:", err)
		ctl.Error(consts.OPERATION_FAILED, "Book Not found")
	}

	ctl.Success(web.M{"Book": result})
}

// Borrow
//
//	@Title			Borrow book
//	@Description	Borrow book
//	@Success		200			{object}	web.M
//	@router			/Borrow [get]
func (ctl *BookController) Borrow() {
	req := dto.ReqBorrow{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[BookController][Borrow] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[BookController][Borrow]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(http.StatusBadRequest)
	}

	db := utility.NewDB()
	tx, err := db.Begin()
	if err != nil {
		logs.Error(err)
	}
	defer tx.Commit()

	loan := models.Loaner{}
	errCode, err := loan.BorrowBook(tx, req)
	if err != nil || errCode != 0 {
		if errCode == 0 {
			errCode = consts.OPERATION_FAILED
		}
		logs.Error("[BookController][Borrow]Something goes wrong:", err)
		ctl.Error(consts.OPERATION_FAILED, "Borrow book error")
	}

	ctl.Success("Success")
}

// Extend
//
//	@Title			Borrow book
//	@Description	Borrow book
//	@Success		200			{object}	web.M
//	@router			/Extend [get]
func (ctl *BookController) Extend() {
	req := dto.ReqExtend{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[BookController][Extend] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[BookController][Extend]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	db := utility.NewDB()
	tx, err := db.Begin()
	if err != nil {
		logs.Error(err)
	}
	defer tx.Commit()

	loan := models.Loaner{}
	errCode, err := loan.ExtendBook(db, tx, req)
	if err != nil || errCode != 0 {
		if errCode == 0 {
			errCode = consts.OPERATION_FAILED
		}
		logs.Error("[BookController][Extend] Something goes wrong:", err)
		ctl.Error(consts.OPERATION_FAILED, "Extend book error")
	}

	ctl.Success("Success")
}

// Return
//
//	@Title			Borrow book
//	@Description	Borrow book
//	@Success		200			{object}	web.M
//	@router			/Return [get]
func (ctl *BookController) Return() {
	req := dto.ReqReturn{}
	if err := ctl.ParseForm(&req); err != nil {
		logs.Error("[BookController][Return] Parse Form Error", err)
		ctl.Error(consts.FAILED_REQUEST)
	}
	if err := validation.ValidateRequest(&req); err != nil {
		logs.Error("[BookController][Return]FormValidate fail, req : %+v, error: %+v", req, err)
		ctl.Error(consts.PARAM_ERROR)
	}

	db := utility.NewDB()
	tx, err := db.Begin()
	if err != nil {
		logs.Error(err)
	}
	defer tx.Commit()

	loan := models.Loaner{}
	errCode, err := loan.ReturnBook(db, tx, req)
	if err != nil || errCode != 0 {
		if errCode == 0 {
			errCode = consts.OPERATION_FAILED
		}
		logs.Error("[BookController][Borrow]Something goes wrong:", err)
		ctl.Error(consts.OPERATION_FAILED, "Borrow book error")
	}

	ctl.Success("Success")
}
