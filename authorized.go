package auth

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Permission is what an authorized peer may do against a served resource.
type Permission int

const (
	// PermRead allows read-only access.
	PermRead Permission = iota
	// PermPropose allows PermRead plus proposing new content without being
	// able to move a ref/pointer directly — narrower than PermWrite.
	PermPropose
	// PermWrite allows PermPropose plus writing refs/pointers directly.
	PermWrite
)

func (p Permission) String() string {
	switch p {
	case PermRead:
		return "read"
	case PermPropose:
		return "propose"
	case PermWrite:
		return "write"
	default:
		return "unknown"
	}
}

func ParsePermission(s string) (Permission, error) {
	switch s {
	case "read":
		return PermRead, nil
	case "propose":
		return PermPropose, nil
	case "write":
		return PermWrite, nil
	default:
		return 0, fmt.Errorf("auth: unknown permission %q (want read, propose, or write)", s)
	}
}

// AuthorizedPeers is a server-side allowlist: fingerprint -> permission.
type AuthorizedPeers map[string]Permission

// LoadAuthorizedPeers parses an authorized-peers file: one
// "<fingerprint> <permission>" pair per line, shaped like
// ~/.ssh/authorized_keys — blank lines and lines starting with # are
// ignored. A missing file is not an error; it just means no peer is
// authorized yet.
func LoadAuthorizedPeers(path string) (AuthorizedPeers, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return AuthorizedPeers{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := AuthorizedPeers{}
	scanner := bufio.NewScanner(f)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("%s:%d: expected \"<fingerprint> <permission>\", got %q", path, lineNo, line)
		}
		perm, err := ParsePermission(fields[1])
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, lineNo, err)
		}
		out[fields[0]] = perm
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Allows reports whether fingerprint is authorized for at least need.
func (a AuthorizedPeers) Allows(fingerprint string, need Permission) bool {
	got, ok := a[fingerprint]
	return ok && got >= need
}
