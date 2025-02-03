package tests

import (
	"bytes"
	initilize "e-library/initialize"
	"e-library/validation"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	booking "e-library/controllers/booking"

	"github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

func init() {
	initilize.InitDB()
	validation.Init()
	web.BConfig.CopyRequestBody = true
}

// Setup Beego router
func setupRouter() {
	bookController := &booking.BookController{}

	web.Router("/All", bookController, "get:All")
	web.Router("/Book", bookController, "get:Book")
	web.Router("/Borrow", bookController, "post:Borrow")
	web.Router("/Extend", bookController, "post:Extend")
	web.Router("/Return", bookController, "post:Return")
}

// Test Get All Books API
func TestAllBooksAPI(t *testing.T) {
	setupRouter()
	req := httptest.NewRequest("GET", "http://localhost:3000/All", nil)
	req.Header.Set("Content-Type", "application/form-data")
	w := httptest.NewRecorder()

	web.BeeApp.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
}

// Test Get Specific Book API
func TestGetBookAPI(t *testing.T) {
	setupRouter()

	// Define query parameters directly in the URL
	req := httptest.NewRequest("GET", "http://localhost:3000/Book?BookTitle=Hello", nil)

	// Set the content type if needed, though for GET requests it may not be necessary
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Set the correct Content-Type

	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
}

// Test Borrow Book API
func TestBorrowBookAPI(t *testing.T) {
	setupRouter()

	// Prepare the form data
	reqData := url.Values{}
	reqData.Set("BookTitle", "Hello")
	reqData.Set("Borrower", "John Doe")

	// Create the request with form-encoded data
	req := httptest.NewRequest("POST", "http://localhost:3000/Borrow", bytes.NewBufferString(reqData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Set the correct Content-Type

	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
}

// Test Extend Borrow Period API
func TestExtendBookAPI(t *testing.T) {
	setupRouter()

	// Prepare the form data
	reqData := url.Values{}
	reqData.Set("BookTitle", "Hello")
	reqData.Set("Borrower", "John Doe")

	// Create the request with form-encoded data
	req := httptest.NewRequest("POST", "http://localhost:3000/Extend", bytes.NewBufferString(reqData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Set the correct Content-Type

	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
}

// Test Return Book API
func TestReturnBookAPI(t *testing.T) {
	setupRouter()

	// Prepare the form data
	reqData := url.Values{}
	reqData.Set("BookTitle", "Hello")
	reqData.Set("Borrower", "John Doe")
	reqData.Set("Value", "1")

	// Create the request with form-encoded data
	req := httptest.NewRequest("POST", "http://localhost:3000/Return", bytes.NewBufferString(reqData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Set the correct Content-Type

	w := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
}
