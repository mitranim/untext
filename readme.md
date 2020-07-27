## Overview

"Missing feature" of the Go packages `encoding` and `fmt`: unmarshal arbitrary
text into an arbitrary value. Counterpart to the marshaling functionality of
`fmt.Sprintf("%v")`.

## Docs

See the full documentation at https://godoc.org/github.com/mitranim/untext.

## Example

```go
var num int64
err := untext.UnmarshalString(`10`, &num)

var inst time.Time
err = untext.UnmarshalString(`0001-02-03T04:05:06Z`, &inst)
```

## License

https://en.wikipedia.org/wiki/WTFPL

## Misc

I'm receptive to suggestions. If this library _almost_ satisfies you but needs changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
