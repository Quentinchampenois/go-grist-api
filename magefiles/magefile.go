//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/quentinchampenois/go-grist-api/examples"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

const (
	BinaryName = "grist-api"
)

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", BinaryName, "./examples/orgs.go")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	fmt.Println("Installing ...")
	return os.Rename(fmt.Sprintf("./%s", BinaryName), fmt.Sprintf("/usr/bin/%s", BinaryName))
}

func Run() error {
	fmt.Fprintln(os.Stdout, "Running ...")
	os.Stdout.Sync() // Ensures the output is flushed
	cmd := exec.Command("go", "run", "./examples/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Orgs() error {
	fmt.Fprintln(os.Stdout, "Running ...")
	os.Stdout.Sync()

	examples.OrgsExample()
	return nil
}

func Workspaces() error {
	fmt.Fprintln(os.Stdout, "Running ...")
	os.Stdout.Sync()

	examples.WorkspacesExample()
	return nil
}

func Docs() error {
	fmt.Fprintln(os.Stdout, "Running ...")
	os.Stdout.Sync()

	examples.DocsExample()
	return nil
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll(BinaryName)
}
