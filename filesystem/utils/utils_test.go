package filesystem

import (
	"testing"
)

func TestExtractURL(t *testing.T) {
	p := ExtractURL("file:///some/place/foo.txt")
	if p.Scheme != "file" || p.Path != "/some/place/foo.txt" {
		t.Fatalf("Expected 'file' and '/some/place/foo.txt'. Got %s and %s", p.Scheme, p.Path)
	}

	p = ExtractURL("/some/place/foo.txt")
	if p.Scheme != "file" || p.Path != "/some/place/foo.txt" {
		t.Fatalf("Expected 'file' and '/some/place/foo.txt'. Got %s and %s", p.Scheme, p.Path)
	}

	p = ExtractURL("http://www.onegeek.com.au/some/place/foo.txt")
	if p.Scheme != "http" || p.Path != "/some/place/foo.txt" || p.Host != "www.onegeek.com.au" {
		t.Fatalf("Expected 'http' and '/some/place/foo.txt'. Got %s and %s", p.Scheme, p.Path)
	}
}

func TestLinuxPath_BasicWindowsPath(t *testing.T) {
	p := LinuxPath(`c:\foo\bar\baz.txt`)
	if p != "/foo/bar/baz.txt" {
		t.Fatalf("Expected /foo/bar/baz.txt but got %s", p)
	}
}

func TestLinuxPath_NoWindowsPrefix(t *testing.T) {
	p := LinuxPath(`\foo\bar\baz.txt`)
	if p != "/foo/bar/baz.txt" {
		t.Fatalf("Expected /foo/bar/baz.txt but got %s", p)
	}
}

func TestLinuxPath_LinuxPath(t *testing.T) {
	p := LinuxPath(`/foo/bar/baz.txt`)
	if p != "/foo/bar/baz.txt" {
		t.Fatalf("Expected /foo/bar/baz.txt but got %s", p)
	}
}
