// main.go

/*
This program is a simple file manager that allows you to get files from a directory and copy them to another directory.
During this process, it will create subdirectories in the target directory based on the creation date of the files.
*/

package main

import (
	"fmt"
	"strings"

	"example.xom/go-demo-1/fileManager"
)

func main() {

	var sourceDir string
	fmt.Print("Bitte geben Sie das Quellverzeichnis ein: ")
	fmt.Scanln(&sourceDir)

	//wenn kein Verzeichnis angegeben wurde, dann nutze das aktuelle Verzeichnis
	if sourceDir == "" {
		sourceDir = "."
	}

	// Get files from the directory
	files := fileManager.GetFilesFromDir(sourceDir)

	// List all found files
	fmt.Println("Alle gefundenen Dateien:")
	files.ListFiles()

	//Frage ob die Dateien verschoben werden sollen
	var moveFiles string
	fmt.Print("Möchten Sie die Dateien verschieben? (j/n): ")
	fmt.Scanln(&moveFiles)

	if strings.ToLower(moveFiles) == "j" {
		// Move files to the target directory
		files.MoveFilesToDir()
		fmt.Println("\nAlle Dateien wurden erfolgreich verschoben.")
	} else {
		fmt.Println("\nVerschieben der Dateien wurde abgebrochen.")
	}
}
