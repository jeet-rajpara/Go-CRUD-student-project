package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Student struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Subject  string `json:"subject"`
	Devision string `json:"devision"`
	Standard string `json:"standard"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS students_data (id SERIAL PRIMARY KEY, name TEXT, subject TEXT, devision TEXT, standard TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// create router
	router := mux.NewRouter()
	router.HandleFunc("/students", getStudentsDetail(db)).Methods("GET")
	router.HandleFunc("/student/{id}", getStudentDetail(db)).Methods("GET")
	router.HandleFunc("/students", addStudentDetail(db)).Methods("POST")
	router.HandleFunc("/students/subject", getStudentsbySubject(db)).Methods("GET")

	log.Println("Successfully server started on port 8005")
	log.Fatal(http.ListenAndServe(":8005", jsonContentTypeMiddleware(router)))

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// get all students details
func getStudentsDetail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT * FROM students_data")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		students := []Student{}
		for rows.Next() {
			var s Student
			if err := rows.Scan(&s.ID, &s.Name, &s.Subject, &s.Devision, &s.Standard); err != nil {
				log.Fatal(err)
			}
			students = append(students, s)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(students)
	}
}

// get student deatil

func getStudentDetail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]

		var s Student
		err := db.QueryRow("SELECT * FROM students_data WHERE id = $1", id).Scan(&s.ID, &s.Name, &s.Subject, &s.Devision, &s.Standard)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(s)
	}
}

//add student detail

func addStudentDetail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var s Student
		json.NewDecoder(r.Body).Decode(&s)
		// fmt.Println(s)

		err := db.QueryRow("INSERT INTO students_data (name, subject, devision, standard) VALUES ($1, $2, $3, $4) RETURNING id", s.Name, s.Subject, s.Devision, s.Standard).Scan(&s.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(s)
	}
}

// get all students based on subject

func getStudentsbySubject(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		subject := r.FormValue("subject")
		rows, err := db.Query("SELECT * FROM students_data WHERE subject = $1", subject)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		students := []Student{}
		for rows.Next() {
			var s Student
			if err := rows.Scan(&s.ID, &s.Name, &s.Subject, &s.Devision, &s.Standard); err != nil {
				log.Fatal(err)
			}
			students = append(students, s)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(students)
	}
}
