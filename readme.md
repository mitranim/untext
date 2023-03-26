**Moved to https://github.com/mitranim/gg**. This repo is usable but frozen.

## Overview

"Missing feature" of the Go packages `encoding` and `fmt`: unmarshal arbitrary
text into an arbitrary value. Counterpart to the marshaling functionality of
`fmt.Sprintf("%v")`.

## Docs

See the full documentation at https://godoc.org/github.com/mitranim/untext.

## Example

```go
var num int64
err := untext.Parse(`10`, &num)

var inst time.Time
err = untext.Parse(`0001-02-03T04:05:06Z`, &inst)
```

## Changelog

### v0.1.3

Breaking: use terms "unmarshal" for `[]byte` inputs and "parse" for `string` inputs. This conforms to the standard library naming conventions.

## License

https://unlicense.org

## Misc

I'm receptive to suggestions. If this library _almost_ satisfies you but needs changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
