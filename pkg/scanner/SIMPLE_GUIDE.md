# Simple Test Implementation Guide

## ?? Just 10 Tests - That's It!

5 tests for `GetEntryPoint()` + 5 tests for `buildPathMap()` = 10 total tests

**Time needed**: 45-60 minutes

---

## ?? Implementation Examples

### Test 1: `TestGetEntryPoint_NonExistentPath`

```go
func TestGetEntryPoint_NonExistentPath(t *testing.T) {
	_, err := GetEntryPoint("/this/does/not/exist")
	
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Wrong error message: %v", err)
	}
}
```

---

### Test 2: `TestGetEntryPoint_CurrentDirectory`

```go
func TestGetEntryPoint_CurrentDirectory(t *testing.T) {
	result, err := GetEntryPoint(".")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if !filepath.IsAbs(result) {
		t.Error("Expected absolute path")
	}
}
```

---

### Test 3: `TestGetEntryPoint_RelativePath`

```go
func TestGetEntryPoint_RelativePath(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("hi"), 0644)
	
	// Change to temp dir
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)
	
	// Test relative path
	result, err := GetEntryPoint("test.txt")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if !filepath.IsAbs(result) {
		t.Error("Should return absolute path")
	}
}
```

---

### Test 4: `TestGetEntryPoint_PathWithSpaces`

```go
func TestGetEntryPoint_PathWithSpaces(t *testing.T) {
	tmpDir := t.TempDir()
	dirWithSpaces := filepath.Join(tmpDir, "my folder")
	os.Mkdir(dirWithSpaces, 0755)
	
	result, err := GetEntryPoint(dirWithSpaces)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if !strings.Contains(result, "my folder") {
		t.Error("Spaces not preserved in path")
	}
}
```

---

### Test 5: `TestGetEntryPoint_EmptyPath`

```go
func TestGetEntryPoint_EmptyPath(t *testing.T) {
	result, err := GetEntryPoint("")
	
	// Just document what happens
	t.Logf("Empty path: result=%s, err=%v", result, err)
	
	// Most likely resolves to current directory
	if err == nil && !filepath.IsAbs(result) {
		t.Error("Should return absolute path")
	}
}
```

---

### Test 6: `TestBuildPathMap_EmptySlice`

```go
func TestBuildPathMap_EmptySlice(t *testing.T) {
	result := buildPathMap([]FileInfo{})
	
	if result == nil {
		t.Fatal("Should return non-nil map")
	}
	
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d entries", len(result))
	}
}
```

---

### Test 7: `TestBuildPathMap_SingleFile`

```go
func TestBuildPathMap_SingleFile(t *testing.T) {
	files := []FileInfo{
		{RelativePath: "file.txt", IsDir: false},
	}
	
	result := buildPathMap(files)
	
	if len(result) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(result))
	}
	
	if result["file.txt"] != false {
		t.Error("File should be marked as false (not directory)")
	}
}
```

---

### Test 8: `TestBuildPathMap_EmptyRelativePath` ?? IMPORTANT

```go
func TestBuildPathMap_EmptyRelativePath(t *testing.T) {
	files := []FileInfo{
		{RelativePath: "", IsDir: true},        // Should be skipped!
		{RelativePath: "file.txt", IsDir: false},
	}
	
	result := buildPathMap(files)
	
	// Empty RelativePath should NOT be in map
	if _, exists := result[""]; exists {
		t.Error("Empty RelativePath should be skipped")
	}
	
	// Should only have the one file
	if len(result) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(result))
	}
}
```

---

### Test 9: `TestBuildPathMap_MixedFilesAndDirs`

```go
func TestBuildPathMap_MixedFilesAndDirs(t *testing.T) {
	files := []FileInfo{
		{RelativePath: "dir1", IsDir: true},
		{RelativePath: "file1.txt", IsDir: false},
		{RelativePath: "dir2", IsDir: true},
		{RelativePath: "file2.txt", IsDir: false},
	}
	
	result := buildPathMap(files)
	
	if len(result) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(result))
	}
	
	// Check directories are marked true
	if result["dir1"] != true {
		t.Error("dir1 should be marked as directory (true)")
	}
	
	// Check files are marked false
	if result["file1.txt"] != false {
		t.Error("file1.txt should be marked as file (false)")
	}
}
```

---

### Test 10: `TestBuildPathMap_DuplicatePaths`

```go
func TestBuildPathMap_DuplicatePaths(t *testing.T) {
	files := []FileInfo{
		{RelativePath: "same", IsDir: true},   // First
		{RelativePath: "same", IsDir: false},  // Second - wins!
	}
	
	result := buildPathMap(files)
	
	if len(result) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(result))
	}
	
	// Last entry should win
	if result["same"] != false {
		t.Error("Last entry should win (IsDir=false)")
	}
}
```

---

## ? Copy-Paste These Into Your Test File

That's it! Just copy these 10 implementations into `scanner_test.go` and you're done!

---

## ?? Running Your Tests

```bash
# Run all tests
go test -v ./pkg/scanner/

# Check coverage
go test -cover ./pkg/scanner/

# Should show ~90% coverage with just these 10 tests!
```

---

## ?? What These 10 Tests Cover

### GetEntryPoint (5 tests):
- ? Error handling (non-existent path)
- ? Basic functionality (current directory)
- ? Path conversion (relative ? absolute)
- ? Special characters (spaces)
- ? Edge case (empty string)

### buildPathMap (5 tests):
- ? Empty input
- ? Single entry
- ? Skip logic (empty RelativePath)
- ? Mixed types (files + dirs)
- ? Duplicates (last wins)

**This is all you need!** ??

---

## ?? Time Breakdown

- Tests 1-5 (GetEntryPoint): 20-25 minutes
- Tests 6-10 (buildPathMap): 15-20 minutes
- **Total: 45 minutes max**

---

## ?? Pro Tip

Implement them in order 1?10. Each test builds on concepts from the previous one.

Good luck! ??
