package models

import (
	"e-library/consts"
	"e-library/utility"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type RentingBook struct {
	CommStruct
	Title           string `orm:"unique;description(Book Title)"`
	AvailableCopies int    `orm:"description(Available unit of the book)"`
	MaxCopies       int    `orm:"description(Max unit of the book)"`
}

type BookDetails struct {
	Title           string
	AvailableCopies int
}

func (book *RentingBook) SetUpdateTime() {
	book.UpdateTime = uint64(time.Now().Unix())
}

func (book *RentingBook) SetCreateTime() {
	book.CreateTime = uint64(time.Now().Unix())
}

func init() {
	orm.RegisterModel(new(RentingBook))
}

func (book *RentingBook) TableName() string {
	return "book"
}

func (book *RentingBook) GetAll() (bookList []BookDetails, errcode int64, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	db := utility.NewDB()

	qb.Select("*")
	qb.From(book.TableName())
	qb.Where("1=1")

	sql := qb.String()
	_, err = db.Raw(sql).QueryRows(&bookList)
	if err != nil {
		logs.Error("[RentingBook][GetAll] Query error:", sql, err)
		errcode = consts.DB_GET_FAILED
		return

	}

	return
}

func (book *RentingBook) GetBook(title string) (bookList RentingBook, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	db := utility.NewDB()

	qb.Select("*")
	qb.From(book.TableName())
	qb.Where("1=1")
	var args []interface{}

	qb.And("title = ?")
	args = append(args, title)

	sql := qb.String()
	err = db.Raw(sql).SetArgs(args).QueryRow(&bookList)
	if err != nil {
		logs.Error("[RentingBook][GetBook] Query error:", sql, err)
		return
	}

	return
}
