package goastap

import "testing"

func TestSolver_VerifyAstapBinary(t *testing.T) {
	solver, err := NewSolver("astap/astap_cli.exe")
	if err != nil {
		t.Fatalf("NewSolver failed: %v", err)
	}
	solver.Solve("testdata/test.fits", false)
}

func TestSolver_SolveDirectory(t *testing.T) {
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
}
