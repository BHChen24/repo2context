package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// =============================================================================
// Tests for GetEntryPoint()
// =============================================================================

func TestGetEntryPoint_NonExistentPath(t *testing.T) {
	// Expected: Should return error

	// Given
	nonExistentPath := "/this/does/not/exist"

	// When
	result, err := GetEntryPoint(nonExistentPath)

	// Then
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != "" {
		t.Fatalf("Expected empty path, got %s", result)
	}
}

func TestGetEntryPoint_CurrentDirectory(t *testing.T) {
	// Expected: Should return absolute path

	// Given
	currentDirectory := "."

	// When
	result, err := GetEntryPoint(currentDirectory)

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty path, got empty")
	}

	if !filepath.IsAbs(result) {
		t.Error("Expected absolute path, got relative")
	}
}

func TestGetEntryPoint_RelativePath(t *testing.T) {
	// Expected: Should convert to absolute path

	// Given
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hi"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if chdirErr := os.Chdir(oldWd); chdirErr != nil {
			t.Fatalf("Failed to restore working directory: %v", chdirErr)
		}
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// When
	result, err := GetEntryPoint(testFile)

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !filepath.IsAbs(result) {
		t.Error("Expected absolute path, got relative")
	}

	if result == "" {
		t.Fatal("Expected non-empty path, got empty")
	}
}

func TestGetEntryPoint_PathWithSpaces(t *testing.T) {
	// Expected: Should handle spaces correctly

	// Given
	tempDir := t.TempDir()
	dirWithSpaces := filepath.Join(tempDir, "my folder")
	if err := os.Mkdir(dirWithSpaces, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// When
	result, err := GetEntryPoint(dirWithSpaces)

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(result, "my folder") {
		t.Errorf("Spaces not preserved in path: %s", result)
	}
}

func TestGetEntryPoint_EmptyPath(t *testing.T) {
	// Expected: Document what happens

	// Given
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to load working directory: %v", err)
	}
	// When
	result, err := GetEntryPoint("")

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !filepath.IsAbs(result) {
		t.Error("Expected absolute path, got relative")
	}

	if result != cwd {
		t.Errorf("Expected current working directory, got %s", result)
	}

}

// =============================================================================
// Tests for buildPathMap()
// =============================================================================

func TestBuildPathMap_EmptySlice(t *testing.T) {
	// Expected: Should return empty map

	// Given
	emptySlice := []FileInfo{}

	// When
	result := buildPathMap(emptySlice)

	// Then
	if result == nil {
		t.Fatal("Expected non-nil map, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d entries", len(result))
	}

}

func TestBuildPathMap_SingleFile(t *testing.T) {
	// Expected: Map should have 1 entry, IsDir=false

	// Given
	files := []FileInfo{
		{RelativePath: "test.txt", IsDir: false},
	}

	// When
	result := buildPathMap(files)

	// Then

	if result == nil {
		t.Fatal("Expected non-nil map, got nil")
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(result))
	}

	val, ok := result["test.txt"]
	if !ok {
		t.Fatalf("Expected \"test.txt\" key to exist")
	}

	if val {
		t.Fatal("Expected file to be marked as false (not directory), got true")
	}

}

func TestBuildPathMap_EmptyRelativePath(t *testing.T) {
	// Expected: Empty path should be SKIPPED

	// Given
	files := []FileInfo{
		{RelativePath: "", IsDir: true},
		{RelativePath: "test.txt", IsDir: false},
	}

	// When
	result := buildPathMap(files)

	// Then
	if result == nil {
		t.Fatal("Expected non-nil map, got nil")
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(result))
	}

	val, ok := result["test.txt"]
	if !ok {
		t.Fatalf("Expected \"test.txt\" key to exist")
	}

	if val {
		t.Fatal("Expected file to be marked as false (not directory), got true")
	}

}

func TestBuildPathMap_MixedFilesAndDirs(t *testing.T) {
	// Expected: Map should correctly mark files (false) and dirs (true)

	// Given
	files := []FileInfo{
		{RelativePath: "dir1", IsDir: true},
		{RelativePath: "file1.txt", IsDir: false},
		{RelativePath: "dir2", IsDir: true},
		{RelativePath: "file2.txt", IsDir: false},
	}

	// When
	result := buildPathMap(files)

	// Then
	if result == nil {
		t.Fatal("Expected non-nil map, got nil")
	}

	if len(result) != 4 {
		t.Fatalf("Expected 4 entries, got %d", len(result))
	}

	cases := map[string]bool{
		"dir1":      true,
		"file1.txt": false,
		"dir2":      true,
		"file2.txt": false,
	}

	for path, expected := range cases {

		val, ok := result[path]
		if !ok {
			t.Fatalf("Expected \"%s\" key to exist", path)
		}

		if val != expected {
			t.Fatalf("Expected %v for \"%s\", got %v", expected, path, val)
		}
	}

}

func TestBuildPathMap_DuplicatePaths(t *testing.T) {
	// Expected: Last entry wins (map overwrites)

	// Given
	files := []FileInfo{
		{RelativePath: "same", IsDir: true},
		{RelativePath: "same", IsDir: false},
	}

	// When
	result := buildPathMap(files)

	// Then
	if result == nil {
		t.Fatal("Expected non-nil map, got nil")
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(result))
	}

	val, ok := result["same"]
	if !ok {
		t.Fatalf("Expected \"same\" key to exist")
	}

	if val {
		t.Fatal("Expected last entry to win (IsDir=false)")
	}
}

// =============================================================================
// Tests for GenerateDirectoryTree()
// =============================================================================

func TestGenerateDirectoryTree_EmptyFiles(t *testing.T) {
	// Expected: Empty string for no files

	// Given
	files := []FileInfo{}
	rootPath := "/dummy"

	// When
	result := generateDirectoryTree(files, rootPath)

	// Then
	expected := ""
	if result != expected {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestGenerateDirectoryTree_SingleFile(t *testing.T) {
	// Expected: Just the file name

	// Given
	files := []FileInfo{
		{RelativePath: "file.txt", IsDir: false, TokenCount: 0},
	}
	rootPath := "/dummy"

	// When
	result := generateDirectoryTree(files, rootPath)

	// Then
	expected := "file.txt\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGenerateDirectoryTree_SingleDirectory(t *testing.T) {
	// Expected: Directory with trailing slash

	// Given
	files := []FileInfo{
		{RelativePath: "dir", IsDir: true, TokenCount: 0},
	}
	rootPath := "/dummy"

	// When
	result := generateDirectoryTree(files, rootPath)

	// Then
	expected := "dir/\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGenerateDirectoryTree_FileWithTokens(t *testing.T) {
	// Expected: File with token count displayed

	// Given
	files := []FileInfo{
		{RelativePath: "file.txt", IsDir: false, TokenCount: 5},
	}
	rootPath := "/dummy"

	// When
	result := generateDirectoryTree(files, rootPath)

	// Then
	expected := "file.txt (5 tokens)\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGenerateDirectoryTree_FileWithZeroTokens(t *testing.T) {
	// Expected: File without token count (since 0)

	// Given
	files := []FileInfo{
		{RelativePath: "file.txt", IsDir: false, TokenCount: 0},
	}
	rootPath := "/dummy"

	// When
	result := generateDirectoryTree(files, rootPath)

	// Then
	expected := "file.txt\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGenerateDirectoryTree_MixedFilesAndDirs(t *testing.T) {
	// Expected: Directories end with /, files show tokens if >0

	// Given
	files := []FileInfo{
		{RelativePath: "dir", IsDir: true, TokenCount: 0},
		{RelativePath: "file1.txt", IsDir: false, TokenCount: 10},
		{RelativePath: "file2.txt", IsDir: false, TokenCount: 0},
	}
	rootPath := "/dummy"

	// When
	result := generateDirectoryTree(files, rootPath)

	// Then
	expected := "dir/\nfile1.txt (10 tokens)\nfile2.txt\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGenerateDirectoryTree_SpecialCharactersInNames(t *testing.T) {
	// Expected: Handle special characters in file/dir names

	// Given
	files := []FileInfo{
		{RelativePath: "file-with-dash.txt", IsDir: false, TokenCount: 0},
		{RelativePath: "dir_with_underscore", IsDir: true, TokenCount: 0},
	}
	rootPath := "/dummy"

	// When
	result := generateDirectoryTree(files, rootPath)

	// Then
	expected := "dir_with_underscore/\nfile-with-dash.txt\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGenerateDirectoryTree_Sorting(t *testing.T) {
	// Expected: Paths sorted alphabetically

	// Given
	files := []FileInfo{
		{RelativePath: "z.txt", IsDir: false, TokenCount: 0},
		{RelativePath: "a.txt", IsDir: false, TokenCount: 0},
		{RelativePath: "m", IsDir: true, TokenCount: 0},
	}
	rootPath := "/dummy"

	// When
	result := generateDirectoryTree(files, rootPath)

	// Then
	expected := "a.txt\nm/\nz.txt\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
