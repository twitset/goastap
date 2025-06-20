package goastap

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Solver struct {
	astapBinaryPath string
}

func NewSolver(astapBinaryPath string) (solver *Solver, err error) {
	solver = &Solver{astapBinaryPath: astapBinaryPath}
	err = solver.verifyBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to verify ASTAP binary: %w", err)
	}
	return
}

func (s *Solver) Solve(imagePath string, disableBackups bool) (err error) {

	if !disableBackups {
		err = createBackupFile(imagePath)
		if err != nil {
			return
		}
	}

	args := []string{
		"-f", imagePath,
		"-update",
		"-wcs",
	}

	// directly invoke the ASTAP binary
	cmd := exec.Command(s.astapBinaryPath, args...)
	if output, err := cmd.CombinedOutput(); err != nil {

		if !disableBackups {
			restoreBackupFile(imagePath)
		}
		return errors.New(string(output))
	}
	return
}

func (s *Solver) SolveDirectory(directoryPath string, disableBackups bool) (notSolvedFilePaths map[string]string, err error) {
	notSolvedFilePaths = make(map[string]string)
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return notSolvedFilePaths, fmt.Errorf("failed to read directory %s: %w", directoryPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories and non-FITS files
		}

		if filepath.Ext(file.Name()) != ".fits" && filepath.Ext(file.Name()) != ".fit" {
			continue // Skip non-FITS files
		}

		imagePath := directoryPath + "/" + file.Name()
		err = s.Solve(imagePath, disableBackups)
		if err != nil {
			notSolvedFilePaths[imagePath] = err.Error()
			continue
		}
	}

	return
}

// verifyBinary checks if the ASTAP binary exists at the specified path.
func (s *Solver) verifyBinary() error {
	if s.astapBinaryPath == "" {
		return fmt.Errorf("astap binary path is not set")
	}

	if _, err := os.Stat(s.astapBinaryPath); os.IsNotExist(err) {
		return fmt.Errorf("astap binary not found at path: %s", s.astapBinaryPath)
	}
	return nil
}

func createBackupFile(filePath string) (err error) {

	backupPath := filePath + ".bak"

	// Open original for reading
	src, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer src.Close()

	// Create backup file
	dst, err := os.Create(backupPath)
	if err != nil {
		return
	}
	defer dst.Close()

	// Copy contents
	if _, errCopy := io.Copy(dst, src); err != nil {
		return errCopy
	}
	return
}

func restoreBackupFile(originalPath string) {
	backupPath := originalPath + ".bak"

	// 1) Does backup exist?
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		fmt.Printf("No backup found for %s, skipping restore\n", originalPath)
		return
	}

	// 2) If an original exists, delete it so Rename can overwrite
	if _, err := os.Stat(originalPath); err == nil {
		if err := os.Remove(originalPath); err != nil {
			fmt.Printf("Error removing existing original %s: %v\n", originalPath, err)
			return
		}
	}

	// 3) Rename .bak → original
	if err := os.Rename(backupPath, originalPath); err != nil {
		fmt.Printf("Error restoring backup: %v\n", err)
		return
	}

	// 4) (Optional) Cleanup: remove leftover .bak if any
	if err := os.Remove(backupPath); err == nil {
		fmt.Printf("Removed backup file: %s\n", backupPath)
	}
}
