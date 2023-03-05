package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Movie struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Movie
}

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/golang?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func getAll(w http.ResponseWriter, r *http.Request) {
	var (
		response Response
		movie    Movie
		movies   []Movie
	)

	db := Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description, rating, image, created_at, updated_at FROM movie")
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
	}

	if err == sql.ErrNoRows {
		response.Status = "Berhasil"
		response.Message = "Tidak ada data"
		response.Data = []Movie{}
	} else {
		for rows.Next() {
			if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Rating, &movie.Image, &movie.CreatedAt, &movie.UpdatedAt); err != nil {
				log.Fatal(err.Error())
			} else {
				movies = append(movies, movie)
			}
		}

		response.Status = "Berhasil"
		response.Message = "Ambil seluruh data Movie"
		response.Data = movies
	}

	fmt.Println("Endpoint Hit: get all Movies")
	json.NewEncoder(w).Encode(response)
}

func getById(w http.ResponseWriter, r *http.Request) {
	var (
		response Response
		movie    Movie
	)

	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	err := db.QueryRow("SELECT id, title, description, rating, image, created_at, updated_at FROM movie WHERE id = ?", id).Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Rating, &movie.Image, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
	}

	if err == sql.ErrNoRows {
		response.Status = "Berhasil"
		response.Message = "Tidak ada data"
		response.Data = []Movie{}
	} else {
		response.Status = "Berhasil"
		response.Message = "Ambil data Movie berdasarkan ID"
		response.Data = []Movie{movie}
	}

	fmt.Println("Endpoint Hit: get by Movie ID")
	json.NewEncoder(w).Encode(response)
}

func create(w http.ResponseWriter, r *http.Request) {
	var (
		response Response
		movie    Movie
	)

	db := Connect()
	defer db.Close()

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &movie)

	message := validateInput(movie, "insert")

	if message != "" {
		response.Status = "Gagal"
		response.Message = message
		response.Data = []Movie{}
	} else {
		_, err := db.Exec("INSERT INTO movie (title, description, rating, image, created_at, updated_at) VALUES (?,?,?,?,?,?,?)", movie.Title, movie.Description, movie.Rating, movie.Image, movie.CreatedAt.Format("2006-01-02 15:04:05"), movie.UpdatedAt.Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Print(err)
		}

		response.Status = "Berhasil"
		response.Message = "Buat data Movie berhasil"
		response.Data = []Movie{movie}
	}

	fmt.Println("Endpoint Hit: create Movie")
	json.NewEncoder(w).Encode(response)
}

func update(w http.ResponseWriter, r *http.Request) {
	var (
		response Response
		movie    Movie
	)

	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	exist := validateById(id, db)
	if exist {
		reqBody, _ := io.ReadAll(r.Body)
		json.Unmarshal(reqBody, &movie)

		message := validateInput(movie, id)

		if message != "" {
			response.Status = "Gagal"
			response.Message = message
			response.Data = []Movie{}
		} else {
			_, err := db.Exec("UPDATE movie SET title = ?, description = ?, rating = ?, image = ?, created_at = ?, updated_at = ? WHERE id = ?", movie.Title, movie.Description, movie.Rating, movie.Image, movie.CreatedAt, movie.UpdatedAt, id)
			if err != nil {
				log.Print(err)
			}

			response.Status = "Berhasil"
			response.Message = "Ubah data Movie berhasil"
			response.Data = []Movie{movie}
		}
	} else {
		response.Status = "Berhasil"
		response.Message = "Tidak ada data"
		response.Data = []Movie{}
	}

	fmt.Println("Endpoint Hit: update Movie")
	json.NewEncoder(w).Encode(response)
}

func delete(w http.ResponseWriter, r *http.Request) {
	var (
		response Response
	)

	vars := mux.Vars(r)
	id := vars["id"]

	db := Connect()
	defer db.Close()

	exist := validateById(id, db)
	if exist {
		_, err := db.Exec("DELETE FROM movie WHERE id = ?", id)
		if err != nil {
			log.Print(err)
		}

		response.Status = "Berhasil"
		response.Message = "Hapus data Movie berhasil"
		response.Data = []Movie{}
	} else {
		response.Status = "Berhasil"
		response.Message = "Tidak ada data"
		response.Data = []Movie{}
	}

	fmt.Println("Endpoint Hit: delete Movie")
	json.NewEncoder(w).Encode(response)
}

func validateById(id string, db *sql.DB) bool {
	var (
		movie Movie
	)

	err := db.QueryRow("SELECT id, title, description, rating, image, created_at, updated_at FROM movie WHERE id = ?", id).Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Rating, &movie.Image, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		return false
	}

	if err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}

func validateInput(movie Movie, id string) string {

	if id != "insert" {
		idNumber, err := strconv.Atoi(id)
		if err != nil {
			log.Print(err)
		}

		if idNumber <= 0 {
			return "ID tidak boleh kosong!"
		}
	} else if movie.Title == "" {
		return "Title tidak boleh kosong!"
	} else if movie.CreatedAt.IsZero() {
		return "Created At tidak boleh kosong!"
	} else if movie.UpdatedAt.IsZero() {
		return "Updated At tidak boleh kosong!"
	}

	return ""
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/Movies", getAll).Methods("GET")
	myRouter.HandleFunc("/Movies", create).Methods("POST")
	myRouter.HandleFunc("/Movies/{id}", update).Methods("PATCH")
	myRouter.HandleFunc("/Movies/{id}", delete).Methods("DELETE")
	myRouter.HandleFunc("/Movies/{id}", getById)
	log.Fatal(http.ListenAndServe(":1234", myRouter))
}

func main() {
	handleRequest()
}
