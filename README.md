# (temporary?) fork of go-webdav

To use this fork, run:
```
go mod edit -replace=github.com/emersion/go-webdav=github.com/oliverpool/go-webdav@main
go mod tidy
```

[![Go Reference](https://pkg.go.dev/badge/github.com/emersion/go-webdav.svg)](https://pkg.go.dev/github.com/emersion/go-webdav)

A Go library for [WebDAV], [CalDAV] and [CardDAV].

## License

MIT

[WebDAV]: https://tools.ietf.org/html/rfc4918
[CalDAV]: https://tools.ietf.org/html/rfc4791
[CardDAV]: https://tools.ietf.org/html/rfc6352
