package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllData(t *testing.T) {
	req, err := http.NewRequest("GET", "/Movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAll)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler yang dikembalikan salah status code: didapatkan %v diinginkan %v", status, http.StatusOK)
	}

	fmt.Println(rr.Body.String())

	expected := `{"status":"Berhasil","message":"Ambil seluruh data Movie","Data":[{"id":1,"title":"Pengabdi Setan 2 Comunion","description":"dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai s","rating":7,"image":"","created_at":"2023-03-06T00:13:11Z","updated_at":"2023-03-06T00:13:11Z"},{"id":2,"title":"Pengabdi Setan","description":"","rating":8,"image":"","created_at":"2023-03-06T00:13:29Z","updated_at":"2023-03-06T00:13:29Z"}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler mengembalikan body yang tidak diharapkan: didapatkan %v diinginkan %v", rr.Body.String(), expected)
	}
}

func TestGetDataById(t *testing.T) {
	req, err := http.NewRequest("GET", "/Movies/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getById)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler yang dikembalikan salah: didapatkan %v diinginkan %v", status, http.StatusOK)
	}

	expected := `{"id":1,"title":"Pengabdi Setan 2 Comunion","description":"dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai s","rating":7,"image":"","created_at":"2023-03-06T00:13:11Z","updated_at":"2023-03-06T00:13:11Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler mengembalikan body yang tidak diharapkan: didapatkan %v diinginkan %v", rr.Body.String(), expected)
	}
}

func TestGeTDataByIdNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/Movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "123")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getById)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler yang dikembalikan salah: didapatkan %v diinginkan %v", status, http.StatusBadRequest)
	}
}

func TestCreateMovie(t *testing.T) {
	var jsonStr = []byte(`{
		"title": "Antman",
		"description": "Manusia Semut",
		"rating": "",
		"image": "",
		"created_at": "2022-08-01T10:56:31Z",
		"updated_at": "2022-08-01T10:56:31Z"
	}`)

	req, err := http.NewRequest("POST", "/Movies", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler yang dikembalikan salah status code: didapatkan %v diinginkan %v", status, http.StatusOK)
	}

	expected := `{"status":"Berhasil","message":"Buat data Movie berhasil","Data":[{"id":3,"title":"Antman","description":"Manusia Semut","rating":0,"image":"","created_at":"2022-08-01T10:56:31Z","updated_at":"2022-08-01T10:56:31Z"}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler yang dikembalikan tidak sesuai harapan body: didapatkan %v diinginkan %v", rr.Body.String(), expected)
	}
}

func TestUpdateMovie(t *testing.T) {
	var jsonStr = []byte(`{
		"id": 3,
		"title": "Antman",
		"description": "Manusia Semut",
		"rating": "",
		"image": "",
		"created_at": "2022-08-01T10:56:31Z",
		"updated_at": 0
	}`)

	req, err := http.NewRequest("PUT", "/Movies", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(update)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler yang dikembalikan salah status code: didapatkan %v diinginkan %v", status, http.StatusOK)
	}

	expected := `{
		"id": 3,
		"title": "Antman",
		"description": "Manusia Semut",
		"rating": "",
		"image": "",
		"created_at": "2022-08-01T10:56:31Z",
		"updated_at": 0
	}`
	if rr.Body.String() != expected {
		t.Errorf("handler yang dikembalikan tidak sesuai harapan body: didapatkan %v diinginkan %v", rr.Body.String(), expected)
	}
}

func TestDeleteMovie(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/Movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "4")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(delete)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler yang dikembalikan salah status code: didapatkan %v diinginkan %v", status, http.StatusOK)
	}

	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler yang dikembalikan tidak sesuai harapan body: didapatkan %v diinginkan %v", rr.Body.String(), expected)
	}
}
