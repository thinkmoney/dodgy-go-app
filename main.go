// / This file is solely to test pipeline tools and should NEVER be called.
package main

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"unsafe"
)

// const `unused` is unused (unused)
const unused = "i am unused"

// CALL THIS FUNCTION AT YOUR OWN PERIL...
// this is all dodgy code to test pipeline code scanners
func main() {
	// func `main` is unused (unused)

	// Setting up err so that it can be referred to easily.
	err := errors.New("")
	fmt.Println(err)

	// Setting up variable to do some dodgyness to
	a := 1
	// ineffectual assignment to a (ineffassign)
	a = 1
	// ineffectual assignment to a (ineffassign)
	a = 2
	// ineffectual assignment to a (ineffassign)
	a = 3
	fmt.Println(a)

	// S1038: should use fmt.Printf instead of fmt.Println(fmt.Sprintf(...)) (but don't forget the newline) (gosimple)
	// missing params to string concat
	// fmt.Println(fmt.Sprintf("%d,%d,%d", a))

	// G101: Potential hardcoded credentials (gosec)
	// gosec uses entropy to try and work out if the value of the hardcoded password looks like a password. The level of entropy can be configured for this rule
	username := "admin"
	password := "aasFFhkjdajhklaIAmAPassword123456789!"
	fmt.Printf("%s:%s", username, password)

	// SA2000: should call wg.Add(1) before starting the goroutine to avoid a race (staticcheck)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
	}()
	wg.Done()

	// G402: TLS MinVersion too low. (gosec)
	tlsConfig := &tls.Config{}
	fmt.Println(tlsConfig)

	// G204: Subprocess launched with a potential tainted input or cmd arguments (gosec)
	cmd := exec.Command("a", "build", "-o", filepath.Join("tmp", "httpserver"), "github.com/docker/docker/contrib/httpserver")
	cmd.Env = []string{"test"}

	// this value of err is never used (SA4006)go-staticcheck
	// ineffectual assignment to err (ineffassign)
	db, err := sql.Open("sqlite3", ":memory:")
	// G202: SQL string concatenation (gosec)
	_, err = db.Query("SELECT * FROM foo where name = " + password)
	// Error return value of `db.Close` is not checked (errcheck)
	db.Close()

	// misusing unsafe block
	intArray := [...]int{1, 2}
	fmt.Printf("\nintArray: %v\n", intArray)
	intPtr := &intArray[0]
	fmt.Printf("\nintPtr=%p, *intPtr=%d.\n", intPtr, *intPtr)
	//G103: Use of unsafe calls should be audited (gosec)
	addressHolder := uintptr(unsafe.Pointer(intPtr)) + unsafe.Sizeof(intArray[0])
	//G103: Use of unsafe calls should be audited (gosec)
	intPtr = (*int)(unsafe.Pointer(addressHolder))

	// infinite loop
	for {
		fmt.Println("I am an infinite loop")
	}

	// unreachable: unreachable code (govet)
	if true {
		fmt.Println("I am reachable")
	} else {
		fmt.Println("I am unreachable")
	}

}
