package goastap

import "testing"

func TestSolver_VerifyAstapBinary(t *testing.T) {
	solver := NewSolver("astap/astap_cli.exe")
	err := solver.VerifyAstapBinary()
	if err != nil {
		t.Errorf("VerifyAstapBinary failed: %v", err)
	} else {
		t.Logf("VerifyAstapBinary succeeded")
	}

	solver.Solve("testdata/not_solved.fits", false)
}

func TestSolver_SolveDirectory(t *testing.T) {
	solver := NewSolver("astap/astap_cli.exe")
	err := solver.SolveDirectory("testdata", false)
	if err != nil {
		t.Errorf("SolveDirectory failed: %v", err)
	} else {
		t.Logf("SolveDirectory succeeded")
	}
}
