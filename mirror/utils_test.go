package mirror

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestGetCustomMirrorDir(t *testing.T) {
	root := "/tmp"
	os.Setenv("MIRROR_HOME", root)
	baseDir := GetMirrorDir()

	if strings.Index(baseDir, root) != 0 {
		t.Fatalf("expected base dir with prefix %s; received %s", root, baseDir)
	}
	os.Setenv("MIRROR_HOME", "")
}

func TestGetMirrorDir(t *testing.T) {
	homeDir := GetHomeDir()
	baseDir := GetMirrorDir()

	if strings.Index(baseDir, homeDir) != 0 {
		t.Fatalf("expected base dir with prefix %s; received %s", homeDir, baseDir)
	}
}

func TestGetCertDir(t *testing.T) {
	root := "/tmp"
	os.Setenv("MIRROR_HOME", root)
	clientDir := GetCertDir()

	if strings.Index(clientDir, root) != 0 {
		t.Fatalf("expected machine client cert dir with prefix %s; received %s", root, clientDir)
	}

	path, filename := path.Split(clientDir)
	if strings.Index(path, root) != 0 {
		t.Fatalf("expected base path of %s; received %s", root, path)
	}
	if filename != "certs" {
		t.Fatalf("expected machine client dir \"certs\"; received %s", filename)
	}
	os.Setenv("MIRROR_HOME", "")
}

func TestGetUsername(t *testing.T) {
	currentUser := "unknown"
	switch runtime.GOOS {
	case "darwin", "linux":
		currentUser = os.Getenv("USER")
	case "windows":
		currentUser = os.Getenv("USERNAME")
	}

	username := GetUsername()
	if username != currentUser {
		t.Fatalf("expected username %s; received %s", currentUser, username)
	}
}

func TestRetryable(t *testing.T) {
	count := 0
	retryMe := func() error {
		t.Logf("RetryMe, attempt number %d", count)
		if count == 2 {
			return nil
		}
		count++
		return errors.New(fmt.Sprintf("Still waiting %d more times...", 2-count))
	}
	retryableSleep = 50 * time.Millisecond
	timeout := 155 * time.Millisecond
	err := Retryable(retryMe, timeout)
	if err != nil {
		t.Fatalf("should not have error retrying function: %s", err.Error())
	}

	count = 0
	timeout = 10 * time.Millisecond
	err = Retryable(retryMe, timeout)
	if err == nil {
		t.Fatalf("should have error retrying funuction")
	}
}
