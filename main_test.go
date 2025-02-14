package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/quack337/goLib/fs"
	"github.com/stretchr/testify/assert"
)

func createFile(t *testing.T, path string, size int, data byte) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	a := make([]byte, size)
	for i := 0; i < size; i++ {
		a[i] = data
	}
	_, err = f.Write(a)
	if err != nil {
		t.Fatal(err)
	}
}

func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "golang_test_")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func createDir(t *testing.T, path string) {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

func checkDstDir(t *testing.T, dir string) {
	files, dirs, err := fs.GetEntries(dir + "/dst")
	if err != nil { t.Fatal(err) }
	assert.True(t, len(files) == 2, fmt.Sprintf("files: %d", len(files)))
	assert.True(t, len(dirs) == 1, fmt.Sprintf("dirs: %d", len(dirs)))

	files, dirs, err = fs.GetEntries(dir + "/dst/sub")
	if err != nil { t.Fatal(err) }
	assert.True(t, len(files) == 2, fmt.Sprintf("files: %d", len(files)))
	assert.True(t, len(dirs) == 0, fmt.Sprintf("dirs: %d", len(dirs)))

	s, err := fs.ReadAllText(dir + "/dst/a.txt")
	if err != nil { t.Fatal(err) }
    assert.True(t, s == "aaaaaaaa", s)

	s, err = fs.ReadAllText(dir + "/dst/b.txt")
	if err != nil { t.Fatal(err) }
    assert.True(t, s == "bbbbbbbbbb", s)

	s, err = fs.ReadAllText(dir + "/dst/sub/c.txt")
	if err != nil { t.Fatal(err) }
    assert.True(t, s == "cccccccccccc", s)

	s, err = fs.ReadAllText(dir + "/dst/sub/d.txt")
	if err != nil { t.Fatal(err) }
    assert.True(t, s == "dddddddddddddd", s)
}

func TestBackupDir(t *testing.T) {
    dir := createTempDir(t)
	defer os.RemoveAll(dir)

    createDir(t, dir + "/src")
    createFile(t, dir + "/src/a.txt", 8, 'a')
    createFile(t, dir + "/src/b.txt", 10, 'b')
    createDir(t, dir + "/src/sub")
    createFile(t, dir + "/src/sub/c.txt", 12, 'c')
    createFile(t, dir + "/src/sub/d.txt", 14, 'd')

    BackupDir(dir + "/src", dir + "/dst")
	checkDstDir(t, dir)

	// delete dst/b.txt
	err := os.Remove(dir + "/dst/b.txt")
	if err != nil { t.Fatal(err) }

	files, dirs, err := fs.GetEntries(dir + "/dst")
	if err != nil { t.Fatal(err) }
	assert.True(t, len(files) == 1)
	assert.True(t, len(dirs) == 1)
	assert.True(t, files[0].Name == "a.txt")

    BackupDir(dir + "/src", dir + "/dst")
	checkDstDir(t, dir)

	// delete dst/sub
	err = os.RemoveAll(dir + "/dst/sub")
	if err != nil { t.Fatal(err) }

	files, dirs, err = fs.GetEntries(dir + "/dst")
	if err != nil { t.Fatal(err) }
	assert.True(t, len(files) == 2)
	assert.True(t, len(dirs) == 0)
	assert.True(t, files[0].Name == "a.txt")
	assert.True(t, files[1].Name == "b.txt")

	BackupDir(dir + "/src", dir + "/dst")
	checkDstDir(t, dir)

	// create dst/c.txt
    createFile(t, dir + "/dst/c.txt", 20, 'c')
	s, err := fs.ReadAllText(dir + "/dst/c.txt")
	if err != nil { t.Fatal(err) }
    assert.True(t, s == "cccccccccccccccccccc")

	BackupDir(dir + "/src", dir + "/dst")
	checkDstDir(t, dir)

	// create dst/sub2/
    createDir(t, dir + "/dst/sub2")
    createFile(t, dir + "/dst/sub2/c.txt", 20, 'c')
	s, err = fs.ReadAllText(dir + "/dst/sub2/c.txt")
	if err != nil { t.Fatal(err) }
    assert.True(t, s == "cccccccccccccccccccc")

	BackupDir(dir + "/src", dir + "/dst")
	checkDstDir(t, dir)
}
