# Changelog

All notable changes to this project are documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project intends to follow [Semantic Versioning](https://semver.org/)
once a first tagged release is cut.

## [Unreleased]

### Added

- Initial extraction of `9vcs`'s `identity` package into its own
  module, generalized so it isn't 9vcs-specific: `Identity`,
  `Fingerprint`/`FingerprintOf`, `ConfigDir` (now `~/.config/9`),
  `Load`, `KnownPeers`/TOFU pinning, `AuthorizedPeers`/`Permission`,
  and TLS 1.3 fingerprint-pinned `ServerTLSConfig`/`ClientTLSConfig`.
  One identity, one trust decision, shared across every 9-family
  program (9vcs, 9sh, ...) instead of a copy per project.
- Legacy migration: `Load()` copies an existing 9vcs-local identity
  at `~/.config/9vcs/identity.{key,cert}` forward to the new
  `~/.config/9` path on first run, preserving the fingerprint (and
  every peer's existing pin of it) across the path change.
