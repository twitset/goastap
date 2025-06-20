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

## Usage

### Importing the package

```go
import (
    "github.com/twitset/goastap"
)


```

### How to use

Take a lock at goastap_test.go for an example of how to use the package. 


## Backup and Restore Behavior

* **Backup**: Before solving, a copy of the original file is made at `<filename>.bak`.
* **Failure Recovery**: If ASTAP returns an error, the backup will be restored automatically (unless backups have been disabled).
* **Cleanup**: After a successful solve, backups remain in place until manually removed.


## Contributing

Contributions, issues, and feature requests are welcome! Feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
