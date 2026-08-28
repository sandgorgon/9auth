package auth

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// migrateLegacyIdentity copies forward an identity created before this
// package was split out of 9vcs into its own module, when one exists.
// 9vcs installs prior to that split keep their identity at
// ~/.config/9vcs/identity.{key,cert}; if keyPath/certPath (the new
// ~/.config/9 location) are both absent but the legacy files exist, they
// are copied over verbatim, preserving the fingerprint — and every peer's
// existing pin of it — across the path change. A no-op once the new path
// has been populated, and a no-op for an install with no legacy identity
// at all.
func migrateLegacyIdentity(keyPath, certPath string) error {
	if _, err := os.Stat(keyPath); err == nil {
		return nil // already migrated (or a fresh install has already generated one)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	base, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	legacyDir := filepath.Join(base, "9vcs")
	legacyKey := filepath.Join(legacyDir, "identity.key")
	legacyCert := filepath.Join(legacyDir, "identity.cert")

	if _, err := os.Stat(legacyKey); errors.Is(err, os.ErrNotExist) {
		return nil // no legacy identity to migrate
	} else if err != nil {
		return err
	}
	if _, err := os.Stat(legacyCert); err != nil {
		return err
	}

	if err := copyFile(legacyKey, keyPath, 0o600); err != nil {
		return err
	}
	return copyFile(legacyCert, certPath, 0o644)
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
