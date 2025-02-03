package models

import (
	"e-library/utility"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type LoanerDetail struct {
	CommStruct
	NameOfBorrower string `orm:"unique;description(Name of people who borrow the book)"`
}

func (loan *LoanerDetail) SetUpdateTime() {
	loan.UpdateTime = uint64(time.Now().Unix())
}

func (loan *LoanerDetail) SetCreateTime() {
	loan.CreateTime = uint64(time.Now().Unix())
}

func init() {
	orm.RegisterModel(new(LoanerDetail))
}

func (loan *LoanerDetail) TableName() string {
	return "loan_detail"
}

func (loan *LoanerDetail) LoanerData(name string) (err error) {

	db := utility.NewDB()
	loan.NameOfBorrower = name

	err = db.Get(loan, "NameOfBorrower")
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[LoanerDetail][LoanerDetail] Query error:", err)
		return
	}

	if err == orm.ErrNoRows {
		loan.NameOfBorrower = name
		loan.SetCreateTime()

		loan.Id, err = db.Insert(loan)
		if err != nil {
			logs.Error("[LoanerDetail][LoanerDetail] Insert new borrower error:", err)
			return
		}
	}
	return
}
