package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	os.RemoveAll("/app/test-1.db")
	os.RemoveAll("/app/test-2.db")
	os.RemoveAll("/app/test-3.db")

	fmt.Println("⛳ pwd")
	output, _ := exec.Command("pwd").Output()
	fmt.Print(string(output))

	fmt.Println("⛳ head /var/opt/tester/companies.db")
	output, _ = exec.Command("head", "/var/opt/tester/companies.db").Output()
	fmt.Print(string(output))

	measureTime("symlink", "/var/opt/tester/companies.db", "/app/test-1.db", symLinkFile)

	measureTime("realSqlite", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", realSqlite)

	measureTime("realSqlite again", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", realSqlite)

	measureTime("db.Query explain", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQueryExplain)

	// measureTime("db.Query (/var/opt/tester/companies.db)", "/var/opt/tester/companies.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQuery)
	// measureTime("db.Query (/app/test-1.db)", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQuery)

	measureTime("cp", "/var/opt/tester/companies.db", "/app/test-3.db", copyFile)
	measureTime("hardlink", "/var/opt/tester/companies.db", "/app/test-2.db", hardLinkFile)
	measureTime("cp again", "/var/opt/tester/companies.db", "/app/test-3.db", copyFile)

	fmt.Println("⛳ end")
}

func measureTime(operation, src, dst string, fn func(string, string) error) {
	start := time.Now()
	fmt.Printf("⛳ Starting %s\n", operation)

	if err := fn(src, dst); err != nil {
		fmt.Printf("- %s failed: %v\n", operation, err)
	} else {
		fmt.Printf("⏰ %v for %s\n", time.Since(start), operation)
	}
}

func copyFile(src, dst string) error {
	cmd := exec.Command("cp", src, dst)
	return cmd.Run()
}

func hardLinkFile(src, dst string) error {
	return os.Link(src, dst)
}

func symLinkFile(src, dst string) error {
	return os.Symlink(src, dst)
}

func realSqlite(src, query string) error {
	output, err := exec.Command("sqlite3", src, ".eqp full", query).Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func dbQueryExplain(src, query string) error {
	db, err := sql.Open("sqlite", src)
	if err != nil {
		fmt.Printf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	// Execute EXPLAIN QUERY PLAN
	rows, err := db.Query("EXPLAIN QUERY PLAN " + query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var detail, from, to, estimatedRows string
		err := rows.Scan(&detail, &from, &to, &estimatedRows)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("💭 EXPLAIN QUERY PLAN: Detail: %s, From: %s, To: %s, Rows: %s\n",
			detail, from, to, estimatedRows)
	}
	return nil
}

func dbQuery(src, query string) error {
	db, err := sql.Open("sqlite", src)
	if err != nil {
		fmt.Printf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	rows, err := db.Query("EXPLAIN QUERY PLAN " + query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var detail, from, to, estimatedRows string
		err := rows.Scan(&detail, &from, &to, &estimatedRows)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Detail: %s, From: %s, To: %s, Rows: %s\n",
			detail, from, to, estimatedRows)
	}

	expectedValues := []string{}
	rows, err = db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var value1 string
		var value2 string

		if err := rows.Scan(&value1, &value2); err != nil {
			return err
		}

		expectedValues = append(expectedValues, strings.Join([]string{value1, value2}, "|"))
	}

	fmt.Println(expectedValues)

	return nil
}
