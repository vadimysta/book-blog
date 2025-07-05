package handlers

import (
	"html/template"
	"log"
	"net/http"
	"database/sql"

	_ "modernc.org/sqlite"
)

type PageDataBooks struct {
	Title string
	Data []Book
}

type PageDataSearch struct {
	Title string
	Query string
	Data []Book
}

type Book struct {
	ID int
	Name string
	Author string
	Published_date string
	Description string
}

func HandlerBooksPage(w http.ResponseWriter, r *http.Request) {

	s := PageDataBooks{
		Title: "Список книжок",
	}

	data, err := sql.Open("sqlite", "data/books.db")
	if err != nil {
		log.Fatalf("Error: invalid data file.db %v", err)
	}

	defer data.Close()

	rows, err := data.Query(`SELECT id, name, author, published_date, description FROM books`)
	if err != nil {
		log.Fatalf("Error: selects statement: %v", err)
	}

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Author, &book.Published_date, &book.Description); err != nil {
			log.Fatalf("Error: scans statement %v", err)
		}
		s.Data = append(s.Data, book)
	}

    tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatalf("Error: invalid files html templates %v", err)
	}

	if err := tmpl.Execute(w, s); err != nil {
		log.Fatalf("Error: executes templates html file %v", err)
	}
}

func HandlerBooksSearchPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	data, err := sql.Open("sqlite", "data/books.db")
	if err != nil {
		log.Fatalf("Error: opens statement %v", err)
	}

	defer data.Close()

	rows, err := data.Query(`SELECT id, name, author, published_date, description FROM books
	WHERE name LIKE ? OR author LIKE ?`, "%"+q+"%", "%"+q+"%")
	if err != nil {
		log.Fatalf("Error: selects statement %v", err)
	}

	defer rows.Close()

	var result []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Author, &book.Published_date, &book.Description); err != nil {
			log.Fatalf("Error: scans statement %v", err)
			continue
		}

		result = append(result, book)
	}

	s := PageDataSearch{
		Title: "Знайдені книжки",
		Query: q,
		Data: result,
	}

	tmpl, err := template.ParseFiles("template/search.html")
	if err != nil {
		log.Fatalf("Error: invalid files html templates %v", err)
	}

	tmpl.Execute(w, s)
}