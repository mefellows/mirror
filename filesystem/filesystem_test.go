package filesystem

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(0)
}
func Test_MockFileSystem(t *testing.T) {
	d1 := []byte("hello\ngo\n")
	ioutil.WriteFile("/tmp/dat1", d1, 0644)

	mock := &MockFileSystem{}
	mockFile := makeFile(false, File{})
	d1 = []byte("hello\ngo\n")
	path := fmt.Sprintf("%s%s-%d", os.TempDir(), "testmockfilesystem-", time.Now().UnixNano())
	ioutil.WriteFile(path, d1, 0644)
	f, _ := os.Open(path)
	defer f.Close()
	r := bufio.NewReader(f)

	mock.ReadBytes = make([]byte, 9)
	r.Read(mock.ReadBytes)
	bytesRead, _ := mock.Read(mockFile)

	if bytes.Compare(bytesRead, d1) != 0 {
		t.Fatalf("File read should read the same set of bytes back. Expected %s, got %s", d1, bytesRead)
	}
	mock.DirError = errors.New("Directory doesn't exist")
	mock.Dir("foo")
	mock.FileTree(File{})
	mock.Write(File{}, make([]byte, 0), 0644)

}

func makeFile(isDir bool, parent File) File {
	f := File{}
	if isDir {
		f.FileMode = os.ModeDir
	} else {
		f.FileMode = 0644
	}
	f.FileSize = rand.Int63()
	prefix := ""
	if parent.Name() != "" {
		prefix = fmt.Sprintf("%s/", parent.Name())
	}
	f.FileName = fmt.Sprintf("%s%s", prefix, makeRandomWord())
	f.FilePath = fmt.Sprintf("%s%s", prefix, makeRandomWord())
	f.FileModTime = time.Now()
	return f
}

func makeRandomWord() string {
	words := []string{"foo", "bar", "bat", "baz", "crab", "cat", "parsnip", "apple", "futon", "peanut", "torture", "reticent", "glassware", "sad", "genius", "toilet", "pan", "chimpanzee", "etc", "var", "camp", "angry", "cloud", "hairy", "jib", "crazy", "counter", "naughty", "wink"}
	return words[rand.Intn(len(words))]
}
