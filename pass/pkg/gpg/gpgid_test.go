package gpg

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mandu/tools/pass/internal/testhelper"
)

func writeGPGID(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, GPGIDFile), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

// isolateGPG points GNUPGHOME at an empty temp dir so that no test can ever
// read the real user's keyring, and clears PASS_GPG_ID.
func isolateGPG(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	if err := os.Chmod(home, 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GNUPGHOME", home)
	t.Setenv("PASS_GPG_ID", "")
}

func TestReadGPGID(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
		wantErr bool
	}{
		{"single", "ABCD1234\n", []string{"ABCD1234"}, false},
		{"no trailing newline", "ABCD1234", []string{"ABCD1234"}, false},
		{"multiple", "AAAA\nBBBB\n", []string{"AAAA", "BBBB"}, false},
		{"comments and blanks", "# store key\n\n  AAAA  \n\n#BBBB\n", []string{"AAAA"}, false},
		{"crlf", "AAAA\r\nBBBB\r\n", []string{"AAAA", "BBBB"}, false},
		{"empty", "", nil, true},
		{"only comments", "# nothing\n\n", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeGPGID(t, dir, tt.content)
			got, found, err := ReadGPGID(dir)
			if !found {
				t.Fatal("found = false, want true")
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadGPGIDMissing(t *testing.T) {
	got, found, err := ReadGPGID(t.TempDir())
	if found || err != nil || got != nil {
		t.Errorf("got (%v, %v, %v), want (nil, false, nil)", got, found, err)
	}
}

func TestResolveRecipientsGPGIDBeatsEnv(t *testing.T) {
	isolateGPG(t)
	t.Setenv("PASS_GPG_ID", "ENVKEY")
	dir := t.TempDir()
	writeGPGID(t, dir, "# store\nSTOREKEY\n")
	got, err := ResolveRecipients(dir)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"STOREKEY"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestResolveRecipientsEnvWhenNoGPGID(t *testing.T) {
	isolateGPG(t)
	t.Setenv("PASS_GPG_ID", "ENVKEY")
	got, err := ResolveRecipients(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"ENVKEY"}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestResolveRecipientsBrokenGPGIDDoesNotFallThrough(t *testing.T) {
	isolateGPG(t)
	t.Setenv("PASS_GPG_ID", "ENVKEY")
	dir := t.TempDir()
	writeGPGID(t, dir, "\n# nothing\n")
	if got, err := ResolveRecipients(dir); err == nil {
		t.Errorf("expected error, got recipients %v", got)
	}
}

func TestResolveRecipientsNoKeys(t *testing.T) {
	if err := CheckGPG(); err != nil {
		t.Skip(err)
	}
	isolateGPG(t) // empty keyring
	dir := t.TempDir()
	_, err := ResolveRecipients(dir)
	if err == nil {
		t.Fatal("expected error with no keys")
	}
	for _, want := range []string{"no GPG secret key", ".gpg-id", "PASS_GPG_ID"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should mention %q", err, want)
		}
	}
}

func TestResolveRecipientsSoleKeyAndAmbiguity(t *testing.T) {
	env, noPassKey, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
	if err != nil {
		t.Skipf("Failed to set up test environment: %v", err)
	}
	defer env.Cleanup()
	t.Setenv("GNUPGHOME", env.GNUPGHome) // throwaway keyring only
	t.Setenv("PASS_GPG_ID", "")
	store := t.TempDir()

	// Two secret keys, no .gpg-id, no PASS_GPG_ID: ambiguous.
	_, err = ResolveRecipients(store)
	if err == nil {
		t.Fatal("expected ambiguity error with two secret keys")
	}
	for _, want := range []string{"multiple GPG secret keys", noPassKey, withPassKey} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should mention %q", err, want)
		}
	}

	// .gpg-id or PASS_GPG_ID resolves the ambiguity.
	writeGPGID(t, store, noPassKey+"\n")
	if got, err := ResolveRecipients(store); err != nil || !reflect.DeepEqual(got, []string{noPassKey}) {
		t.Errorf(".gpg-id: got (%v, %v)", got, err)
	}
	os.Remove(filepath.Join(store, GPGIDFile))
	t.Setenv("PASS_GPG_ID", withPassKey)
	if got, err := ResolveRecipients(store); err != nil || !reflect.DeepEqual(got, []string{withPassKey}) {
		t.Errorf("env: got (%v, %v)", got, err)
	}
	t.Setenv("PASS_GPG_ID", "")

	// Delete one secret key (in the throwaway keyring): the other is now the sole key.
	// (batch mode requires the full fingerprint)
	fprOut, err := exec.Command("gpg", "--batch", "--list-keys", "--with-colons", withPassKey).Output()
	if err != nil {
		t.Fatal(err)
	}
	var fpr string
	for _, line := range strings.Split(string(fprOut), "\n") {
		if strings.HasPrefix(line, "fpr:") {
			fpr = strings.Split(line, ":")[9]
			break
		}
	}
	if fpr == "" {
		t.Fatal("no fingerprint found in test keyring")
	}
	del := exec.Command("gpg", "--batch", "--yes", "--delete-secret-keys", fpr)
	if out, err := del.CombinedOutput(); err != nil {
		t.Fatalf("failed to delete secret key in test keyring: %v\n%s", err, out)
	}
	got, err := ResolveRecipients(store)
	if err != nil {
		t.Fatalf("sole key: %v", err)
	}
	if !reflect.DeepEqual(got, []string{noPassKey}) {
		t.Errorf("sole key: got %v, want [%s]", got, noPassKey)
	}
}

// TestEncryptUsesGPGIDNotEnv verifies, with throwaway keys, that a secret is
// encrypted to the .gpg-id key even though PASS_GPG_ID names another key.
func TestEncryptUsesGPGIDNotEnv(t *testing.T) {
	env, noPassKey, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
	if err != nil {
		t.Skipf("Failed to set up test environment: %v", err)
	}
	defer env.Cleanup()
	t.Setenv("GNUPGHOME", env.GNUPGHome)
	t.Setenv("PASS_GPG_ID", noPassKey)

	store := t.TempDir()
	writeGPGID(t, store, withPassKey+"\n")
	recipients, err := ResolveRecipients(store)
	if err != nil {
		t.Fatal(err)
	}

	src := filepath.Join(store, "plain.txt")
	if err := os.WriteFile(src, []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(store, "out.gpg")
	opts := BatchGPGOptions("")
	opts.Recipients = recipients
	if err := EncryptFileWithOptions(src, dst, opts); err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	// Only the passphrase-protected key is a recipient.
	got, err := DecryptFileWithOptions(dst, BatchGPGOptions("test-passphrase-123"))
	if err != nil {
		t.Fatalf("decrypt with .gpg-id key failed: %v", err)
	}
	if got != "secret" {
		t.Errorf("got %q, want %q", got, "secret")
	}
}

func TestEncryptAmbiguousKeyringFails(t *testing.T) {
	env, _, _, err := testhelper.SetupTestEnvWithGPGKeys()
	if err != nil {
		t.Skipf("Failed to set up test environment: %v", err)
	}
	defer env.Cleanup()
	t.Setenv("GNUPGHOME", env.GNUPGHome)
	t.Setenv("PASS_GPG_ID", "")

	dir := t.TempDir()
	src := filepath.Join(dir, "plain.txt")
	if err := os.WriteFile(src, []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "out.gpg")
	if err := EncryptFileWithOptions(src, dst, BatchGPGOptions("")); err == nil {
		t.Error("expected error: two secret keys and no recipient")
	}
	if _, err := os.Stat(dst); err == nil {
		t.Error("no output file may be written")
	}
}
