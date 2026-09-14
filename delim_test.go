package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDelimiter(t *testing.T) {
	tests := []struct {
		name     string
		char     string
		width    int
		expected string
	}{
		{
			name:     "standard single character",
			char:     "=",
			width:    10,
			expected: "==========",
		},
		{
			name:     "custom character hyphen",
			char:     "-",
			width:    5,
			expected: "-----",
		},
		{
			name:     "multi-character pattern",
			char:     "=-",
			width:    3,
			expected: "=-=-=-",
		},
		{
			name:     "unicode character",
			char:     "─",
			width:    4,
			expected: "────",
		},
		{
			name:     "zero width",
			char:     "=",
			width:    0,
			expected: "",
		},
		{
			name:     "negative width",
			char:     "=",
			width:    -10,
			expected: "",
		},
		{
			name:     "empty char",
			char:     "",
			width:    20,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDelimiter(tt.char, tt.width)
			if got != tt.expected {
				t.Errorf("buildDelimiter(%q, %d) = %q, want %q", tt.char, tt.width, got, tt.expected)
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	t.Run("explicit width positive", func(t *testing.T) {
		w := getTerminalWidth(42)
		if w != 42 {
			t.Errorf("expected 42, got %d", w)
		}
	})

	t.Run("columns env var fallback", func(t *testing.T) {
		t.Setenv("COLUMNS", "132")
		w := getTerminalWidth(0)
		if w != 132 {
			t.Errorf("expected 132, got %d", w)
		}
	})

	t.Run("default fallback when columns unset", func(t *testing.T) {
		t.Setenv("COLUMNS", "")
		w := getTerminalWidth(0)
		if w <= 0 {
			t.Errorf("expected positive terminal width fallback, got %d", w)
		}
	})
}

func TestRun(t *testing.T) {
	t.Run("default output with width and newline", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{"-w", "10", "-n"}, &stdout, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "==========\n"
		if stdout.String() != expected {
			t.Errorf("stdout = %q, want %q", stdout.String(), expected)
		}
	})

	t.Run("custom char with short flag", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{"-c", "*", "-w", "5"}, &stdout, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "*****"
		if stdout.String() != expected {
			t.Errorf("stdout = %q, want %q", stdout.String(), expected)
		}
	})

	t.Run("custom char with long flag", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{"--char", "#", "--width", "4"}, &stdout, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "####"
		if stdout.String() != expected {
			t.Errorf("stdout = %q, want %q", stdout.String(), expected)
		}
	})

	t.Run("version flag", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{"-v"}, &stdout, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(stdout.String(), "delim") {
			t.Errorf("expected version output to contain 'delim', got %q", stdout.String())
		}
	})

	t.Run("help flag", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{"-h"}, &stdout, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(stdout.String(), "Usage:") {
			t.Errorf("expected help output to contain 'Usage:', got %q", stdout.String())
		}
	})

	t.Run("displayconfig flag", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{"--displayconfig"}, &stdout, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(stdout.String(), "CONFIG: char :") {
			t.Errorf("expected displayconfig to contain 'CONFIG: char :', got %q", stdout.String())
		}
	})

	t.Run("invalid width flag", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{"-w", "-5"}, &stdout, &stderr)
		if err == nil {
			t.Fatalf("expected error for negative width, got nil")
		}
	})

	t.Run("explicit missing config file", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{"--config", "/nonexistent/config/file.yaml"}, &stdout, &stderr)
		if err == nil {
			t.Fatalf("expected error for non-existent explicit config, got nil")
		}
		if !strings.Contains(err.Error(), "config file not found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("load custom config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "custom.yaml")
		if err := os.WriteFile(configFile, []byte("char: '~'\n"), 0644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		var stdout, stderr bytes.Buffer
		err := run([]string{"--config", configFile, "-w", "5"}, &stdout, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stdout.String() != "~~~~~" {
			t.Errorf("stdout = %q, want %q", stdout.String(), "~~~~~")
		}
	})

	t.Run("environment variable DELIM_CHAR", func(t *testing.T) {
		t.Setenv("DELIM_CHAR", "@")
		var stdout, stderr bytes.Buffer
		err := run([]string{"-w", "6"}, &stdout, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stdout.String() != "@@@@@@" {
			t.Errorf("stdout = %q, want %q", stdout.String(), "@@@@@@")
		}
	})
}

func BenchmarkBuildDelimiter_SingleChar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = buildDelimiter("=", 120)
	}
}

func BenchmarkBuildDelimiter_MultiChar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = buildDelimiter("=-", 120)
	}
}

func BenchmarkRun(b *testing.B) {
	args := []string{"-w", "80"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = run(args, io.Discard, io.Discard)
	}
}
