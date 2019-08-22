package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	// Add trace module
	"./trace"
)

type Employee struct {
	Id    int
	Name  string
	City string
}

func CallStep2(writer http.ResponseWriter, request *http.Request) {
	employees := getData(request)
	for _, employee := range employees {
		fmt.Fprintf(writer, "ID: %d\t Name: %s\t City: %s\n", employee.Id, employee.Name, employee.City)
	}
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser   := "root"
	dbPass   := "password"
	dbName   := "people"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getData(r *http.Request) []Employee {


	connectionSpan := trace.TraceDBConnection(r, "Database connection")
	db := dbConn()
	connectionSpan.Finish()


	sql := "SELECT * FROM Employee ORDER BY id DESC"
	sqlSpan := trace.TraceSQLExecution(r, sql, "SQL execution")
	selDB, err := db.Query(sql)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	sqlSpan.Finish()
	defer db.Close()
	return res
}

func main() {
		log.Println("Server started on: http://localhost:8082")
		trace.SetupTracer("step2")
		// Added Tracing Handler
		http.HandleFunc("/", trace.GoSensor.TracingHandler("step1", CallStep2))

		http.ListenAndServe(":8082", nil)
}

