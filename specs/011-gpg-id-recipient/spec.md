# `.gpg-id` Recipient Selection Specification

## Overview

`pass insert` and `pass edit` (which re-encrypts) must encrypt secrets to the
GPG key(s) listed in the `.gpg-id` file in the root of the password store.
The existing `PASS_GPG_ID` environment variable and default-key behaviour
remain as fallbacks for stores that have no `.gpg-id`, but the fallback to
"the default key" is made safe: it is only allowed when the choice is
unambiguous.

## Status

- **Status**: Implemented ✅
- **Created**: 2026-08-10
- **Branch**: `fix/011-gpg-id-recipient`

## Background

`pass-edit/spec.md` defines `PASS_GPG_ID` ("GPG recipient key ID", default
"Default key") and the implementation resolves the recipient as
`opts.Recipient`, then `PASS_GPG_ID`, then `--default-recipient-self`. The
specs never mention `.gpg-id`, which is how the reference `pass` binds a
store to its key(s), so `insert`/`edit` ignore it.

## Problem Statement

1. `.gpg-id` in the store root is ignored; secrets are encrypted to whatever
   key `PASS_GPG_ID`/the default key selects, not the store's key.
2. With no `PASS_GPG_ID`, `--default-recipient-self` silently picks one of
   possibly several secret keys, so the wrong key can be used.

## Goals

- If `<store>/.gpg-id` exists, `insert` and `edit` encrypt to it.
- Otherwise the existing behaviour applies: `PASS_GPG_ID` if set.
- Otherwise use the user's GPG key, but only if there is exactly one usable
  secret key. With none, or more than one, fail with a clear error.

## Non-Goals

- Per-subdirectory `.gpg-id` files: only the store root is read.
- Creating `.gpg-id` / initialising a store (use the reference `pass init`).
- Verifying the keys named in `.gpg-id` up front (gpg reports problems when encrypting).
- Re-encrypting existing secrets when `.gpg-id` changes.

## Design

### Recipient resolution (first match wins)

1. `<store>/.gpg-id`: one key ID, fingerprint or e-mail per line; blank lines
   and `#` comments ignored; whitespace (incl. `\r`) trimmed.
2. `PASS_GPG_ID` environment variable.
3. The only usable secret key in the keyring (`gpg --list-secret-keys
   --with-colons`; keys whose validity is revoked `r`, expired `e` or
   disabled `d` are not counted). Its key ID is passed explicitly as
   `--recipient`.

A `.gpg-id` that exists but is unreadable or contains no recipient is an
error; resolution does **not** fall through to 2/3, since that would silently
encrypt to a different key than the store declares.

`GPGOptions.Recipient` (explicit single recipient, used by tests) keeps its
current meaning and sits above `PASS_GPG_ID` but below an explicit
`GPGOptions.Recipients` (added here, carries the `.gpg-id` result).
`--default-recipient-self` is no longer used.

### Errors

| Condition | Message |
|---|---|
| `.gpg-id` unreadable / no recipients | `pass: invalid .gpg-id: <reason>` |
| no `.gpg-id`, no `PASS_GPG_ID`, no secret key | `pass: no GPG secret key found; create <store>/.gpg-id or set PASS_GPG_ID` |
| no `.gpg-id`, no `PASS_GPG_ID`, several secret keys | `pass: multiple GPG secret keys found (<ids>); cannot tell which to use. Create <store>/.gpg-id or set PASS_GPG_ID` |

`insert` resolves recipients before prompting for the secret; `edit` before
opening the editor. Neither prints anything to stdout on failure.

### API

`pkg/gpg`:

```go
func ReadGPGID(storeDir string) (recipients []string, found bool, err error)
func ResolveRecipients(storeDir string) ([]string, error) // chain above
```

`GPGOptions.Recipients []string` is added. `EncryptFileWithOptions` uses, in
order: `Recipients`, `Recipient`, `PASS_GPG_ID`, the sole secret key; or
returns an error.

`cmd`: `storeEncryptOptions(storeDir)` calls `ResolveRecipients` and returns
options for `insert` and `edit`.

## Testing

No test touches the user's keyring or password store. Tests use throwaway
keys in a temp `GNUPGHOME` (from `testhelper`) or no GPG at all.

- `ReadGPGID`: missing, single, multiple, comments, blanks, CRLF, empty.
- `ResolveRecipients`: `.gpg-id` beats `PASS_GPG_ID`; env used when no
  `.gpg-id`; broken `.gpg-id` errors even with env set; sole key used when it
  is the only one (second key's secret deleted in the temp keyring); two keys
  and no env/`.gpg-id` -> ambiguity error naming both; no keys -> error.
- Encryption goes to the `.gpg-id` key rather than the `PASS_GPG_ID` key
  (decrypt with the former's passphrase).
- Existing tests that set `PASS_GPG_ID` keep passing unchanged.
