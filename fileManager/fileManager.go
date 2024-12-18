package fileManager

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	Filename  string
	SourceDir string
	TargetDir string
}

type Files []File

func New() File {
	return File{"Filename", "SourceDir", "TargetDir"}
}

func (f Files) ListFiles() {
	for _, file := range f {
		fmt.Printf("Filename: %s | Source Directory: %s | Target Directory: %s\n",
			file.Filename,
			file.SourceDir,
			file.TargetDir)
	}
}

func GetFilesFromDir(srcDir string) Files {
	var files Files
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return files
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		file := File{
			Filename:  entry.Name(),
			SourceDir: srcDir,
			TargetDir: filepath.Join(srcDir, fmt.Sprintf("%d-%02d", info.ModTime().Year(), info.ModTime().Month())),
		}
		files = append(files, file)
	}

	return files
}

func (f *Files) MoveFilesToDir() {
	for _, file := range *f {
		srcPath := filepath.Join(file.SourceDir, file.Filename)
		dstPath := filepath.Join(file.TargetDir, file.Filename)

		// Prüfen, ob Quelldatei existiert
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			fmt.Printf("Quelldatei existiert nicht: %s\n", srcPath)
			continue
		}

		// Zielverzeichnis erstellen
		err := os.MkdirAll(file.TargetDir, 0755)
		if err != nil {
			fmt.Printf("Fehler beim Erstellen des Verzeichnisses %s: %v\n", file.TargetDir, err)
			continue
		}

		// Erst versuchen, die Datei direkt zu verschieben
		err = os.Rename(srcPath, dstPath)
		if err != nil {
			// Falls Rename fehlschlägt (z.B. bei verschiedenen Laufwerken),
			// auf Kopieren und Löschen zurückfallen
			err = moveFileByCopy(srcPath, dstPath)
			if err != nil {
				fmt.Printf("Fehler beim Verschieben der Datei von %s nach %s: %v\n", srcPath, dstPath, err)
				continue
			}
		}
	}
}

// moveFileByCopy kopiert die Datei und löscht das Original
func moveFileByCopy(srcPath, dstPath string) error {
	source, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("fehler beim öffnen der quelldatei: %v", err)
	}
	defer source.Close()

	destination, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("fehler beim erstellen der zieldatei: %v", err)
	}
	defer destination.Close()

	srcInfo, err := source.Stat()
	if err != nil {
		return fmt.Errorf("fehler beim lesen der dateiattribute: %v", err)
	}

	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("fehler beim kopieren der datei: %v", err)
	}

	destination.Close()

	if err := os.Chmod(dstPath, srcInfo.Mode()); err != nil {
		return fmt.Errorf("fehler beim setzen der berechtigungen: %v", err)
	}

	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("fehler beim löschen der originaldatei: %v", err)
	}

	return nil
}
