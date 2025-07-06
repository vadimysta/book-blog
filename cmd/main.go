package main

import (
	"book-blog/handlers"
	"log"
	"net/http"
	"database/sql"

	"github.com/gorilla/mux"
)

func initDB() {
	data, err := sql.Open("sqlite", "data/books.db")
	if err != nil {
		log.Fatalf("Error: invalid data file.db %v", err)
	}

	defer data.Close()

	t, err := data.Prepare(`CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		author TEXT NOT NULL,
		published_date TEXT NOT NULL,
		description TEXT
	);`)
	if err != nil {
		log.Fatalf("Error: preparing statement: %v", err)
	}

	if _, err := t.Exec(); err != nil {
		log.Fatalf("Error: executing statement: %v", err)
	}

	var count int
	row := data.QueryRow("SELECT COUNT(*) FROM books")
	if err := row.Scan(&count); err != nil {
		log.Fatalf("Error: counting rows: %v", err)
	}

	if count == 0 {
		_, err := data.Exec(`
		INSERT INTO books(name, author, published_date, description) VALUES
		('Інтернат', 'Сергій Жадан', '2017-09-01', 'Роман про події на сході України, про пошуки людяності посеред війни.'),
		('Танець недоумка', 'Ілларіон Павлюк', '2016-03-15', 'Психологічний трилер про журналіста і серійного вбивцю.'),
		('Солодка Даруся', 'Марія Матіос', '2004-05-10', 'Трагічна історія гуцульської дівчини на тлі історичних подій.'),
		('Планета Полин', 'Оксана Забужко', '2000-10-20', 'Збірка есеїв про українське суспільство, культуру та ідентичність.'),
		('Сміття', 'Дмитро Скочко', '2019-07-12', 'Соціальна драма про людей, що живуть серед смітників.')
		`)
		if err != nil {
			log.Fatalf("Error: inserts statement: %v", err)
		}
	}
}

func main() {

	initDB()

	router := mux.NewRouter()

	router.HandleFunc("/list/books", handlers.HandlerBooksPage).Methods("GET")
	router.HandleFunc("/list/books/search", handlers.HandlerBooksSearchPage).Methods("GET")
	router.HandleFunc("/list/books/add", handlers.HandlerBooksAddPage).Methods("GET")
	router.HandleFunc("/list/books/add", handlers.HandlerBooksSubmitPage).Methods("POST")


	log.Fatal(http.ListenAndServe(":7070", router))
}