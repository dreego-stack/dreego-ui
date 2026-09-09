package main

import (
	"flag"
	"fmt"
	"os"

	ui "github.com/dreego-stack/dreego-ui"
)

func main() {
	install := flag.NewFlagSet("install", flag.ExitOnError)
	destination := install.String("destination", "", "component destination directory")
	if len(os.Args) < 2 || os.Args[1] != "install" {
		fmt.Fprintln(os.Stderr, "usage: dreego-ui install [-destination path]")
		os.Exit(2)
	}
	if err := install.Parse(os.Args[2:]); err != nil {
		os.Exit(2)
	}
	paths, err := ui.InstallComponents(ui.InstallOptions{Destination: *destination})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, path := range paths {
		fmt.Println(path)
	}
}
