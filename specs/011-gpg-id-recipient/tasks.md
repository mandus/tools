# Tasks: `.gpg-id` recipient selection

- [x] Spec (`spec.md`)
- [x] `gpg.ReadGPGID`, `GPGOptions.Recipients`
- [x] `gpg.ResolveRecipients` (.gpg-id > PASS_GPG_ID > sole secret key; ambiguity error)
- [x] `EncryptFileWithOptions` uses the chain; drop `--default-recipient-self`
- [x] `storeEncryptOptions` in `cmd`; use in `insert.go` and `edit.go`
- [x] Unit tests (no GPG) and throwaway-key tests (temp `GNUPGHOME`)
- [x] CHANGELOG, `docs/pass-quick-reference.md`, `specs/pass-edit/spec.md` env table
