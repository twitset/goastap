package goastap

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Solver struct {
	astapBinaryPath string
}

func NewSolver(astapBinaryPath string) *Solver {
	if astapBinaryPath == "" {
		panic("astap binary path cannot be found")
	}
	return &Solver{astapBinaryPath: astapBinaryPath}
}

func (s *Solver) Solve(imagePath string, disableBackups bool) {

	if !disableBackups {
		createBackupFile(imagePath)
	}

	args := []string{
		"-f", imagePath,
		"-update",
		"-wcs",
	}

	// directly invoke the ASTAP binary
	cmd := exec.Command(s.astapBinaryPath, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("ASTAP failed: %v\nOutput: %s\n", err, output)

		// if we fail, we restore the backup file if it was created
		if !disableBackups {
			restoreBackupFile(imagePath)
		}
		return
	}
}

func (s *Solver) SolveDirectory(directoryPath string, disableBackups bool) error {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", directoryPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories and non-FITS files
		}

		if filepath.Ext(file.Name()) != ".fits" && filepath.Ext(file.Name()) != ".fit" {
			continue // Skip non-FITS files
		}

		imagePath := directoryPath + "/" + file.Name()
		s.Solve(imagePath, disableBackups)
	}

	return nil
}

// VerifyAstapBinary checks if the ASTAP binary exists at the specified path.
func (s *Solver) VerifyAstapBinary() error {
	if s.astapBinaryPath == "" {
		return fmt.Errorf("astap binary path is not set")
	}

	if _, err := os.Stat(s.astapBinaryPath); os.IsNotExist(err) {
		return fmt.Errorf("astap binary not found at path: %s", s.astapBinaryPath)
	}
	return nil
}

func createBackupFile(filePath string) {

	backupPath := filePath + ".bak"

	// Open original for reading
	src, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening original file: %v\n", err)
		return
	}
	defer src.Close()

	// Create backup file
	dst, err := os.Create(backupPath)
	if err != nil {
		fmt.Printf("Error creating backup file: %v\n", err)
		return
	}
	defer dst.Close()

	// Copy contents
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Printf("Error copying to backup: %v\n", err)
		return
	}

	fmt.Printf("Backup created at: %s\n", backupPath)
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
		fmt.Printf("Removed existing original: %s\n", originalPath)
	}

	// 3) Rename .bak → original
	if err := os.Rename(backupPath, originalPath); err != nil {
		fmt.Printf("Error restoring backup: %v\n", err)
		return
	}
	fmt.Printf("Restored backup to: %s\n", originalPath)

	// 4) (Optional) Cleanup: remove leftover .bak if any
	if err := os.Remove(backupPath); err == nil {
		fmt.Printf("Removed backup file: %s\n", backupPath)
	}
}
