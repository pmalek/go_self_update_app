# go_self_update_app

## How to build

Adjust the `Version` constant in the `main.go` to match the one in the filename

### *nix and Mac OS

Build with:

```
go build -o server_v<integer_version_number>
```

### Windows

Build with:

```
go build -o server_v<integer_version_number>.exe
```

## Test

### With Docker

Assuming that we have built `server_v1` and `server_v2` we could test it
via docker like so:

```
cd target/linux
docker run --rm -p 8080:8080 -v $(pwd):/app -w /app alpine ./server_v1
```

## Assumptions

* Server binaries are assumed to have the following format:
  * `server_v<integer_version_number>` for \*nix and Mac OS
  * `server_v<integer_version_number>.exe` for Windows
* Version number is embedded in the binary with `Version` constant in `main.go`

## Room for improvement

* Manage the version const and version embedded in filename better. Ideally hold
  version only in one place. For the filesystem update provider we could store
  version information in some metadata on the fs and embed the version information
  in the server binary somehow instead of relying on using filenames.
* Locking: there's some unnecessary locking done in the handler. This is OK for PoC
  but should be rethought when deployed.
* Custom javascript redirect (2s after clicking install new version) should add
  some more visual indication that the update is underway and shouldn't rely on
  constant time redirect but should rather periodically check if server is back up.
