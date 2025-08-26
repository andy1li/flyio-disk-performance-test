package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	os.RemoveAll("/app/test-1.db")
	os.RemoveAll("/app/test-2.db")
	os.RemoveAll("/app/test-3.db")

	// fmt.Println("⛳ pwd")
	// output, _ := exec.Command("pwd").Output()
	// fmt.Print(string(output))

	// fmt.Println("⛳ head /var/opt/tester/companies.db")
	// output, _ = exec.Command("head", "/var/opt/tester/companies.db").Output()
	// fmt.Print(string(output))

	measureTime("symlink", "/var/opt/tester/companies.db", "/app/test-1.db", symLinkFile)

	// file, err := os.OpenFile("/app/test-1.db", os.O_RDWR, 0644)
	// if err != nil {
	// 	fmt.Printf("Failed to open database: %v\n", err)
	// 	return
	// }
	// defer file.Close()

	// measureTimeForReadPage(file, 0)

	// i := 128
	// for i < 252249 {
	// 	for j := 0; j < 128; j++ {
	// 		measureTimeForReadPage(file, i+j)
	// 	}
	// 	i *= 2
	// }

	// measureTime("realSqlite limit 1", "./test-1.db", "SELECT id, name FROM companies LIMIT 1", realSqlite)
	// // measureTime("realSqlite", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", realSqlite)
	// // measureTime("realSqlite again", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", realSqlite)

	measureTime("db.Query limit 1", "./test-1.db", "SELECT id, name FROM companies LIMIT 1", dbQuery)
	measureTime("db.Query index", "/test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQuery)
	// measureTime("db.Query explain", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQueryExplain)
	// // measureTime("db.Query (/var/opt/tester/companies.db)", "/var/opt/tester/companies.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQuery)
	// measureTime("db.Query (./test-1.db)", "./test-1.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQuery)
	// // measureTime("db.Query (./companies.db)", "./companies.db", "SELECT id, name FROM companies WHERE country = 'micronesia'", dbQuery)

	// // measureTime("cp", "/var/opt/tester/companies.db", "/app/test-3.db", copyFile)
	// measureTime("hardlink", "/var/opt/tester/companies.db", "/app/test-2.db", hardLinkFile)
	// measureTime("cp again", "/var/opt/tester/companies.db", "/app/test-3.db", copyFile)

	// fmt.Println("⛳ end")
}

func measureTime(operation, src, dst string, fn func(string, string) error) {
	start := time.Now()
	fmt.Printf("⛳ Starting %s\n", operation)

	if err := fn(src, dst); err != nil {
		fmt.Printf("❌ %s failed: %v\n", operation, err)
	} else {
		fmt.Printf("⏰ %v for %s\n", time.Since(start), operation)
	}
}

func measureTimeForReadPage(file *os.File, pageNumber int) error {
	start := time.Now()

	if err := readPage(file, pageNumber); err != nil {
		fmt.Printf("❌ readPage %d failed: %v\n", pageNumber, err)
	} else {
		fmt.Printf("⏰ %v for reading 📄 %d\n", time.Since(start), pageNumber)
	}

	return nil
}

func copyFile(src, dst string) error {
	cmd := exec.Command("cp", src, dst)
	return cmd.Run()
}

func dbQuery(src, query string) error {
	fmt.Println("⛳ before sql.Open")

	db, err := sql.Open("sqlite", src)

	fmt.Println("⛳ after sql.Open")

	timeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("⏰", time.Now())
		panic(fmt.Sprintf("❌ CodeCrafters internal error: The tester failed to open the test database within %v", timeout))
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

func hardLinkFile(src, dst string) error {
	return os.Link(src, dst)
}

func readPage(file *os.File, pageNumber int) error {
	offset := int64(pageNumber * 4096)

	_, err := file.Seek(offset, 0) // 0 = seek from beginning
	if err != nil {
		return fmt.Errorf("failed to seek to page %d: %v", pageNumber, err)
	}

	buffer := make([]byte, 4096)
	bytesRead, err := io.ReadFull(file, buffer)
	if err != nil {
		return fmt.Errorf("failed to read page %d: %v", pageNumber, err)
	}

	fmt.Printf("📄 %d: read %d bytes. First 16 bytes: %x\n", pageNumber, bytesRead, buffer[:16])
	return nil
}

func realSqlite(src, query string) error {
	output, err := exec.Command("sqlite3", src, ".eqp full", query).Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func symLinkFile(src, dst string) error {
	return os.Symlink(src, dst)
}
