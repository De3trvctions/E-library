package models

import (
	"e-library/consts"
	"e-library/models/dto"
	"e-library/utility"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type Loaner struct {
	CommStruct
	LoanerId         int64  `orm:"description(LoanerId)"`
	LoanDate         uint64 `orm:"description(Date of the loaning)"`
	ReturnDate       uint64 `orm:"description(Date of the book should be returning)"`
	ActualReturnDate uint64 `orm:"description(Actual book returning date)"`
	BookId           int64  `orm:"description(Book ID which the people is loaning)"`
	Returned         bool   `orm:"default(false);description(Has the book been returned?)"`
}

func (loan *Loaner) SetUpdateTime() {
	loan.UpdateTime = uint64(time.Now().Unix())
}

func (loan *Loaner) SetCreateTime() {
	loan.CreateTime = uint64(time.Now().Unix())
}

func init() {
	orm.RegisterModel(new(Loaner))
}

func (loan *Loaner) TableName() string {
	return "loan"
}

func (loan *Loaner) BorrowBook(tx *utility.TxOrm, req dto.ReqBorrow) (errCode int64, err error) {
	book := RentingBook{}
	result, err := book.GetBook(req.BookTitle)
	if result.Id == 0 && result.AvailableCopies == 0 {
		return
	}

	loanerDetail := LoanerDetail{}
	err = loanerDetail.LoanerData(req.Borrower)
	if err != nil {
		logs.Error("[Loaner][BorrowBook] Get borrower error", err)
		return
	}

	result.AvailableCopies = result.AvailableCopies - 1
	_, err = tx.Update(&result, "AvailableCopies")
	if err != nil {
		logs.Error("[Loaner][BorrowBook] Update book error. %+v , error: %+v", result, err)
		return 123, err
	}

	dateTime := time.Now()
	loan.LoanerId = loanerDetail.Id
	loan.LoanDate = uint64(dateTime.Unix())
	loan.ReturnDate = uint64(dateTime.Add(4 * 7 * 24 * time.Hour).Unix())
	loan.BookId = result.Id
	loan.SetCreateTime()

	_, err = tx.Insert(loan)
	if err != nil {
		logs.Error("[Loaner][BorrowBook] Insert load record error", err)
		return 123, err
	}

	return
}

func (loan *Loaner) ExtendBook(db *utility.DB, tx *utility.TxOrm, req dto.ReqExtend) (errCode int64, err error) {
	qbCheckBook, _ := orm.NewQueryBuilder("mysql")
	qbCheckBook.Select("A.*")
	qbCheckBook.From(loan.TableName() + " AS A")
	qbCheckBook.LeftJoin("book AS B").On("B.id = A.book_id")
	qbCheckBook.LeftJoin("loan_detail AS C").On("C.id = A.loaner_id")
	qbCheckBook.Where("1=1")
	var args []interface{}
	var data []Loaner

	qbCheckBook.And("B.title = ?")
	args = append(args, req.BookTitle)

	qbCheckBook.And("C.name_of_borrower = ?")
	args = append(args, req.Borrower)

	qbCheckBook.And("A.return_date >= ?")
	args = append(args, uint64(time.Now().Unix()))

	qbCheckBook.And("A.returned = false")

	_, err = db.Raw(qbCheckBook.String()).SetArgs(args).QueryRows(&data)
	if err != nil {
		logs.Error("[Loaner][ExtendBook] Check extenable book error", err)
		return consts.DB_GET_FAILED, err
	}

	if len(data) > 0 {
		for _, v := range data {
			originalTime := time.Unix(int64(v.ReturnDate), 0)
			newTime := originalTime.Add(3 * 7 * 24 * time.Hour).Unix()

			v.ReturnDate = uint64(newTime)
			v.SetUpdateTime()
			_, err = tx.Update(&v, "ReturnDate", "UpdateTime")
			if err != nil {
				logs.Error("[Loaner][ExtendBook] Extend Book error", v, err)
				return consts.DB_GET_FAILED, err
			}
		}
	}

	return
}

func (loan *Loaner) ReturnBook(db *utility.DB, tx *utility.TxOrm, req dto.ReqReturn) (errCode int64, err error) {
	qbCheckBook, _ := orm.NewQueryBuilder("mysql")
	qbCheckBook.Select("A.*")
	qbCheckBook.From(loan.TableName() + " AS A")
	qbCheckBook.LeftJoin("book AS B").On("B.id = A.book_id")
	qbCheckBook.LeftJoin("loan_detail AS C").On("C.id = A.loaner_id")
	qbCheckBook.Where("1=1")
	var args []interface{}
	var data []Loaner

	qbCheckBook.And("B.title = ?")
	args = append(args, req.BookTitle)

	qbCheckBook.And("C.name_of_borrower = ?")
	args = append(args, req.Borrower)

	qbCheckBook.And("A.returned = false")
	qbCheckBook.Limit(req.Value)

	_, err = db.Raw(qbCheckBook.String()).SetArgs(args).QueryRows(&data)
	if err != nil {
		logs.Error("[Loaner][ReturnBook] Check return book error", err)
		return consts.DB_GET_FAILED, err
	}

	bookMap := make(map[int64]int)
	for _, v := range data {
		if _, exists := bookMap[v.BookId]; !exists {
			bookMap[v.BookId] = 1
		} else {
			bookMap[v.BookId] += 1
		}

		v.ActualReturnDate = uint64(time.Now().Unix())
		v.SetUpdateTime()
		_, err = tx.Update(&v, "ActualReturnDate", "UpdateTime")
		if err != nil {
			logs.Error("[Loaner][ReturnBook] Return Book error", v, err)
			return consts.DB_UPDATE_FAILED, err
		}
	}

	for i, v := range bookMap {
		updateBook := RentingBook{}
		updateBook.Id = i
		err = db.Get(&updateBook)
		if err != nil {
			logs.Error("[Loaner][ReturnBook] Get Book error", v, err)
			return consts.DB_GET_FAILED, err
		}

		if updateBook.AvailableCopies+v > updateBook.MaxCopies {
			logs.Error("[Loaner][ReturnBook] Return book more than maximum value")
			return consts.OPERATION_FAILED, errors.New("Return book more than maximum value")
		}

		updateBook.AvailableCopies = updateBook.AvailableCopies + v

		updateBook.SetUpdateTime()
		_, err = tx.Update(&updateBook, "AvailableCopies", "UpdateTime")
		if err != nil {
			logs.Error("[Loaner][ReturnBook] Update return book available copies error")
			return consts.DB_UPDATE_FAILED, err
		}
	}

	return
}
