# 9auth

[![CI (master)](https://github.com/sandgorgon/9auth/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/sandgorgon/9auth/actions/workflows/ci.yml?query=branch%3Amaster)
[![CI (develop)](https://github.com/sandgorgon/9auth/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/sandgorgon/9auth/actions/workflows/ci.yml?query=branch%3Adevelop)

A pure-Go, zero-dependency identity and peer-trust primitive shared
across every 9-family program (`9vcs`, `9sh`, and others to come):
one per-install Ed25519 identity, one trust decision, instead of a
copy per project.

```
go get github.com/sandgorgon/9auth
```

## What it provides

- `Identity` — a long-lived Ed25519 keypair wrapped in a minimal
  self-signed X.509 certificate, persisted at
  `~/.config/9/identity.{key,cert}`. `Load()` generates one on first
  use.
- `Fingerprint(pub)` / `(*Identity).Fingerprint()` / `FingerprintOf(cert)`
  — a stable SHA-256-of-public-key identifier for out-of-band
  exchange ("what's your fingerprint?").
- `ServerTLSConfig` / `ClientTLSConfig` — TLS 1.3 configs that
  authenticate the peer by exact fingerprint match instead of a CA
  chain. There is no CA in this design.
- `AuthorizedPeers` — a server-side allowlist (`fingerprint ->
  Permission`), file format modeled on `~/.ssh/authorized_keys`.
- `KnownPeers` — a client-side TOFU pin (`address -> fingerprint`),
  file format modeled on `~/.ssh/known_hosts`.

An install that already has an identity at the legacy
`~/.config/9vcs/identity.{key,cert}` path (from before this package
existed as its own module) has it copied forward automatically on
first `Load()` — the fingerprint, and every peer's existing pin of
it, survives the move.

## Quick start

```go
id, err := auth.Load()
if err != nil {
    log.Fatal(err)
}
fmt.Println("this install's fingerprint:", id.Fingerprint())

authorized, _ := auth.LoadAuthorizedPeers(authorizedPeersPath)
cfg := id.ServerTLSConfig(func(fp string) bool {
    return authorized.Allows(fp, auth.PermRead)
})
```

## Design

There is no certificate authority. Trust is established
out-of-band (a human compares fingerprints, or a TOFU prompt records
one on first connect) and enforced purely by exact fingerprint match
during the TLS handshake — `VerifyPeerCertificate` is the only gate,
so an unauthorized peer never reaches the application layer at all.

## Testing

```
go build ./...
go vet ./...
go test -race ./...
```
