package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yinxulai/go-template/internal/extension"
)

func main() {
	installFlag := flag.String("i", "", "Install extension from zip file")
	uninstallFlag := flag.Bool("u", false, "Uninstall extension")
	flag.Parse()

	if *installFlag != "" {
		zipPath, err := filepath.Abs(*installFlag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if err := extension.Install(zipPath); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else if *uninstallFlag {
		if err := extension.Uninstall("NewEngine"); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage:")
		fmt.Println("  Install extension: cei -i <path_to_zip>")
		fmt.Println("  Uninstall extensions: cei -u")
	}
}
