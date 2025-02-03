package dto

type ReqBorrow struct {
	Borrower  string `valid:"Required"`
	BookTitle string `valid:"Required"`
}

type ReqExtend struct {
	Borrower  string `valid:"Required"`
	BookTitle string `valid:"Required"`
}

type ReqGetBook struct {
	BookTitle string `valid:"Required"`
}

type ReqReturn struct {
	Borrower  string `valid:"Required"`
	BookTitle string `valid:"Required"`
	Value     int    `valid:"Required"`
}
