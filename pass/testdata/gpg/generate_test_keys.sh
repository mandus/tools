#!/bin/bash
# Generate test GPG keys for testing
# This script creates test-only GPG keys that can be used in automated tests
# NEVER use these keys for real password storage

set -euo pipefail

# Configuration
TEST_GNUPG_HOME="./test-gnupg-home"
KEY_NAME="Test User"
KEY_EMAIL="test@example.com"
KEY_COMMENT="Test Key"

# Passphrase for test key with passphrase
TEST_PASSPHRASE="test-passphrase-123"

# Clean up existing test home
rm -rf "$TEST_GNUPG_HOME"
mkdir -p "$TEST_GNUPG_HOME"
chmod 700 "$TEST_GNUPG_HOME"

echo "Generating test GPG keys in $TEST_GNUPG_HOME"

# Export current GPG_TTY if set
export GNUPGHOME="$TEST_GNUPG_HOME"

# Disable gpg-agent for batch operations
export GPG_TTY=$(tty 2>/dev/null || echo "")

# Generate a test key WITHOUT passphrase (for basic tests)
echo "Generating test key without passphrase..."
gpg --batch \
    --gen-key <<EOF
Key-Type: RSA
Key-Length: 2048
Subkey-Type: RSA
Subkey-Length: 2048
Name-Real: $KEY_NAME
Name-Email: $KEY_EMAIL
Name-Comment: $KEY_COMMENT (No Passphrase)
Expire-Date: 0
%no-protection
%commit
EOF

# Get the key ID of the generated key
NO_PASSPHRASE_KEYID=$(gpg --list-keys --with-colons | grep '^pub' | head -1 | cut -d: -f5)
echo "Generated key without passphrase: $NO_PASSPHRASE_KEYID"

# Generate a test key WITH passphrase (for passphrase testing)
echo "Generating test key with passphrase..."
gpg --batch \
    --passphrase "$TEST_PASSPHRASE" \
    --gen-key <<EOF
Key-Type: RSA
Key-Length: 2048
Subkey-Type: RSA
Subkey-Length: 2048
Name-Real: $KEY_NAME Passphrase
Name-Email: test-passphrase@example.com
Name-Comment: $KEY_COMMENT (With Passphrase)
Expire-Date: 0
%commit
EOF

# Get the key ID of the passphrase-protected key
WITH_PASSPHRASE_KEYID=$(gpg --list-keys --with-colons | grep '^pub' | tail -1 | cut -d: -f5)
echo "Generated key with passphrase: $WITH_PASSPHRASE_KEYID"

# Export the keys
gpg --export --armor "$NO_PASSPHRASE_KEYID" > no-passphrase-key.pub
gpg --export-secret-keys --armor "$NO_PASSPHRASE_KEYID" > no-passphrase-key.sec

gpg --export --armor "$WITH_PASSPHRASE_KEYID" > with-passphrase-key.pub
gpg --export-secret-keys --armor "$WITH_PASSPHRASE_KEYID" > with-passphrase-key.sec

# Create a config file for batch mode
cat > "$TEST_GNUPG_HOME/gpg.conf" <<EOF
batch
no-tty
yes
EOF

# Create a gpg-agent config that won't prompt
cat > "$TEST_GNUPG_HOME/gpg-agent.conf" <<EOF
pinentry-program /bin/false
EOF

# Output key IDs for Go tests to use
cat > "$TEST_GNUPG_HOME/key-ids.txt" <<EOF
NO_PASSPHRASE_KEYID=$NO_PASSPHRASE_KEYID
WITH_PASSPHRASE_KEYID=$WITH_PASSPHRASE_KEYID
TEST_PASSPHRASE=$TEST_PASSPHRASE
EOF

echo ""
echo "Test keys generated successfully!"
echo "Key without passphrase: $NO_PASSPHRASE_KEYID"
echo "Key with passphrase: $WITH_PASSPHRASE_KEYID"
echo ""
echo "To use these keys in tests:"
echo "  export GNUPGHOME=$(pwd)/$TEST_GNUPG_HOME"
echo ""
