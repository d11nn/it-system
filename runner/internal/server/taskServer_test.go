package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsTestSuccessDetectsGoBuildFailures(t *testing.T) {
	server := &taskServer{}

	testCases := []struct {
		name    string
		output  string
		success bool
	}{
		{
			name:    "pass output",
			output:  "PASS\nok  \ttest\t0.123s\n",
			success: true,
		},
		{
			name:    "tab separated build failure",
			output:  "# github.com/free5gc/amf/internal/gmm\nhandler.go:1: cannot use int32 as *int32\nFAIL\ttest [build failed]\n",
			success: false,
		},
		{
			name:    "space separated build failure",
			output:  "FAIL github.com/free5gc/free5gc/test [build failed]\n",
			success: false,
		},
		{
			name:    "ansi wrapped build failure",
			output:  "\x1b[31mFAIL\ttest [build failed]\x1b[0m\n",
			success: false,
		},
		{
			name:    "go test failure",
			output:  "--- FAIL: TestRegistration (0.01s)\n",
			success: false,
		},
		{
			name:    "exit status failure",
			output:  "exit status 1\n",
			success: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := server.isTestSuccess(tc.output); got != tc.success {
				t.Fatalf("isTestSuccess() = %v, want %v", got, tc.success)
			}
		})
	}
}

func TestBuildLibraryReplaceArg(t *testing.T) {
	got := buildLibraryReplaceArg("openapi", "contributor/openapi", "abc123")
	want := "-replace=github.com/free5gc/openapi=github.com/contributor/openapi@abc123"
	if got != want {
		t.Fatalf("buildLibraryReplaceArg() = %q, want %q", got, want)
	}
}

func TestGoModuleDirs(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{
		root,
		filepath.Join(root, "test"),
		filepath.Join(root, "NFs", "amf"),
		filepath.Join(root, "NFs", "smf"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
			t.Fatalf("write go.mod in %s: %v", dir, err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git", "go.mod"), []byte("module ignored\n"), 0o644); err != nil {
		t.Fatalf("write .git go.mod: %v", err)
	}

	got, err := goModuleDirs(root)
	if err != nil {
		t.Fatalf("goModuleDirs(): %v", err)
	}

	want := []string{
		root,
		filepath.Join(root, "NFs", "amf"),
		filepath.Join(root, "NFs", "smf"),
		filepath.Join(root, "test"),
	}
	if len(got) != len(want) {
		t.Fatalf("goModuleDirs() length = %d, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("goModuleDirs()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
