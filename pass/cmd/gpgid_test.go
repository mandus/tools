package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// These tests avoid invoking gpg: GNUPGHOME points at an empty temp dir.
func isolateGPG(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	if err := os.Chmod(home, 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GNUPGHOME", home)
	t.Setenv("PASS_GPG_ID", "")
}

func TestStoreEncryptOptions(t *testing.T) {
	t.Run(".gpg-id is used and beats PASS_GPG_ID", func(t *testing.T) {
		isolateGPG(t)
		t.Setenv("PASS_GPG_ID", "ENVKEY")
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ".gpg-id"), []byte("# k\nAAAA\nBBBB\n"), 0600); err != nil {
			t.Fatal(err)
		}
		opts, err := storeEncryptOptions(dir)
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"AAAA", "BBBB"}; !reflect.DeepEqual(opts.Recipients, want) {
			t.Errorf("Recipients = %v, want %v", opts.Recipients, want)
		}
	})

	t.Run("PASS_GPG_ID used when no .gpg-id", func(t *testing.T) {
		isolateGPG(t)
		t.Setenv("PASS_GPG_ID", "ENVKEY")
		opts, err := storeEncryptOptions(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"ENVKEY"}; !reflect.DeepEqual(opts.Recipients, want) {
			t.Errorf("Recipients = %v, want %v", opts.Recipients, want)
		}
	})

	t.Run("empty .gpg-id is an error even with PASS_GPG_ID", func(t *testing.T) {
		isolateGPG(t)
		t.Setenv("PASS_GPG_ID", "ENVKEY")
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ".gpg-id"), []byte("\n"), 0600); err != nil {
			t.Fatal(err)
		}
		_, err := storeEncryptOptions(dir)
		if err == nil || !strings.Contains(err.Error(), "invalid .gpg-id") {
			t.Errorf("err = %v, want invalid .gpg-id error", err)
		}
	})
}
