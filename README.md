# ZipIt

A small HTTP service that decrypts a URL, downloads the remote file, wraps it in a ZIP archive, and streams the archive back to the client.

## Requirements

- Go 1.25 or newer
- `ENCRYPTION_KEY`: exactly 32 characters

## Run

```bash
export ENCRYPTION_KEY="0123456789abcdef0123456789abcdef"
export PORT=8080   # optional
go run .
```

## Usage

Encrypt the target URL with the same 32-character key. The payload must be
AES-256-CBC encrypted, with the IV prepended and the result encoded as
base64url without padding. `Generate-String.ps1` is a small helper for
Windows/PowerShell.

Then request:

```text
http://localhost:8080/<encrypted-url>
```

The service downloads the file and returns it as `download.zip`.

## Security

This service fetches arbitrary URLs after decrypting the path. Do **not**
expose it directly to the public internet: it can be abused as an open proxy
or SSRF vector. Run it behind authentication and an allowlist, and only fetch
URLs you trust.

## License

MIT. See [LICENSE](LICENSE).
