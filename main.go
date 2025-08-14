package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	os.RemoveAll("/tmp/test-1.db")
	os.RemoveAll("/tmp/test-2.db")
	os.RemoveAll("/tmp/test-3.db")

	measureTime("symlink", "./companies.db", "/tmp/test-1.db", symLinkFile)
	measureTime("hardlink", "./companies.db", "/tmp/test-2.db", hardLinkFile)
	measureTime("cp", "./companies.db", "/tmp/test-3.db", copyFile)
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
