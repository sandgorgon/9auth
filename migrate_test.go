package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateLegacyIdentityCopiesForward(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	base, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}

	legacyDir := filepath.Join(base, "9vcs")
	if err := os.MkdirAll(legacyDir, 0o700); err != nil {
		t.Fatal(err)
	}
	legacy, err := loadFrom(filepath.Join(legacyDir, "identity.key"), filepath.Join(legacyDir, "identity.cert"))
	if err != nil {
		t.Fatal(err)
	}

	newDir := filepath.Join(base, "9")
	keyPath := filepath.Join(newDir, "identity.key")
	certPath := filepath.Join(newDir, "identity.cert")

	if err := migrateLegacyIdentity(keyPath, certPath); err != nil {
		t.Fatal(err)
	}
	migrated, err := loadFrom(keyPath, certPath)
	if err != nil {
		t.Fatal(err)
	}
	if migrated.Fingerprint() != legacy.Fingerprint() {
		t.Fatalf("migrated fingerprint %s, want %s (legacy)", migrated.Fingerprint(), legacy.Fingerprint())
	}
}

func TestMigrateLegacyIdentityNoLegacyIsNoop(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	base, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}

	keyPath := filepath.Join(base, "9", "identity.key")
	certPath := filepath.Join(base, "9", "identity.cert")
	if err := migrateLegacyIdentity(keyPath, certPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(keyPath); !os.IsNotExist(err) {
		t.Fatalf("expected no key to be created, got err=%v", err)
	}
}

func TestMigrateLegacyIdentityDoesNotOverwriteExisting(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	base, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}

	legacyDir := filepath.Join(base, "9vcs")
	if err := os.MkdirAll(legacyDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := loadFrom(filepath.Join(legacyDir, "identity.key"), filepath.Join(legacyDir, "identity.cert")); err != nil {
		t.Fatal(err)
	}

	newDir := filepath.Join(base, "9")
	if err := os.MkdirAll(newDir, 0o700); err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(newDir, "identity.key")
	certPath := filepath.Join(newDir, "identity.cert")
	current, err := loadFrom(keyPath, certPath) // already-established new-path identity
	if err != nil {
		t.Fatal(err)
	}

	if err := migrateLegacyIdentity(keyPath, certPath); err != nil {
		t.Fatal(err)
	}
	after, err := loadFrom(keyPath, certPath)
	if err != nil {
		t.Fatal(err)
	}
	if after.Fingerprint() != current.Fingerprint() {
		t.Fatalf("migration overwrote an already-established identity: got %s, want %s", after.Fingerprint(), current.Fingerprint())
	}
}
