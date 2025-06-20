# goastap

Go wrapper for invoking the [ASTAP](http://www.hnsky.org/astap.htm) command‑line tool to perform astrometric plate solving and World Coordinate System (WCS) updates on FITS images.

## Features

* Solve a single FITS/FIT file and update its WCS headers by calling the ASTAP binary
* Recursively process all FITS files in a directory
* Automatic backup of originals before solving, with optional backup disable flag
* Verification helper to ensure the ASTAP executable is present

## Prerequisites

* Go 1.18 or higher
* ASTAP CLI executable and its catalog files must be placed together in the same directory (the "ASTAP directory").
* Download both the ASTAP command-line client and at least one compatible star catalog from the official site: [https://www.hnsky.org/astap.htm](https://www.hnsky.org/astap.htm)

## Installation

```bash
go get github.com/twitset/goastap
```

### Importing the package

```go
import (
    "github.com/twitset/goastap"
)


```

## Usage

### Single Image
```go 

solver, err := NewSolver("astap/astap_cli.exe") // default is "astap/astap_cli.exe" if left empty
if err != nil {
	log.Println("Error creating solver:", err)
}
err = solver.Solve("testdata/test.fits", false)
if err != nil {
    log.Println("Error solving image:", err)
} else {
    log.Println("Image solved successfully!")
}

```

### Directory Processing
```go
solver, err := NewSolver("astap/astap_cli.exe")
if err != nil {
    t.Fatalf("NewSolver failed: %v", err)
}
notSolvedPaths, err := solver.SolveDirectory("testdata", false)
if err != nil {
    t.Errorf("SolveDirectory failed: %v", err)
} else {
    t.Logf("SolveDirectory succeeded")
}

if len(notSolvedPaths) > 0 {
    t.Errorf("SolveDirectory found unsolved files: %v", notSolvedPaths)
} else {
    t.Logf("All files solved successfully")
}
````

### Examples

Take a lock at goastap_test.go for an example of how to use the package. 


## Backup and Restore Behavior

* **Backup**: Before solving, a copy of the original file is made at `<filename>.bak`.
* **Failure Recovery**: If ASTAP returns an error, the backup will be restored automatically (unless backups have been disabled).
* **Cleanup**: After a successful solve, backups remain in place until manually removed.

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
