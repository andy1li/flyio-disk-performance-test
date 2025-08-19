package main

import (
	"database/sql"
	"fmt"
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
	measureTime("db.Query (/var/opt/tester/companies.db)", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQuery)
	measureTime("db.Query (/app/test-1.db)", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQuery)

	measureTime("cp", "/var/opt/tester/companies.db", "/app/test-3.db", copyFile)
	measureTime("hardlink", "/var/opt/tester/companies.db", "/app/test-2.db", hardLinkFile)

	fmt.Println("⛳ end")
}

func measureTime(operation, src, dst string, fn func(string, string) error) {
	start := time.Now()
	fmt.Printf("Starting %s\n", operation)

	if err := fn(src, dst); err != nil {
		fmt.Printf("- %s failed: %v\n", operation, err)
	} else {
		fmt.Printf("- %v for %s\n", time.Since(start), operation)
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

func dbQuery(src, query string) error {
	db, err := sql.Open("sqlite", src)
	if err != nil {
		fmt.Printf("Failed to create test database, this is a CodeCrafters error.")
		return err
	}
	defer db.Close()

	expectedValues := []string{}

	rows, err := db.Query(query)
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
