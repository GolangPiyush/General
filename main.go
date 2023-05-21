package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {

	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "db_fmbs_app"
	dbPass := "FmbsAppDbPwd"
	dbName := "db_fmbs"

	// Construct connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Open a connection to the database
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Drop the employee table if it exists
	_, err = db.Exec("DROP TABLE IF EXISTS employees")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Employee table dropped!")

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")

	// Create the employee table
	_, err = db.Exec(`
		CREATE TABLE employees (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Employee table created!")

	// Insert an employee
	employee := Employee{
		ID:   1,
		Name: "John Doe",
	}
	_, err = db.Exec("INSERT INTO employees (id, name) VALUES (?, ?)", employee.ID, employee.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Employee inserted successfully!")

	// Perform a simple query
	rows, err := db.Query("SELECT id, name FROM employees")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/employee", func(ctx *gin.Context) {
		var employees []Employee
		rows, err := db.Query("SELECT id, name FROM employees")
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		defer rows.Close()

		// Iterate over the result set
		for rows.Next() {
			var id int
			var name string
			err := rows.Scan(&id, &name)
			if err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
			employees = append(employees, Employee{ID: id, Name: name})
		}

		ctx.JSON(http.StatusOK, employees)
	})

	r.POST("/employee", func(ctx *gin.Context) {
		var emp Employee
		if err := ctx.BindJSON(&emp); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		_, err := db.Exec("INSERT INTO employees VALUES(?,?)", emp.ID, emp.Name)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ctx.JSON(http.StatusOK, emp)
	})
	r.Run(":8080")
}
