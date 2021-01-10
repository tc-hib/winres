package version

const (
	errInvalidSignature    = "invalid fixed file info signature"
	errInvalidLength       = "invalid length"
	errInvalidLangID       = "invalid language id"
	errUnhandledCodePage   = "unhandled code page"
	errInvalidStringLength = "invalid string length"

	errEmptyKey         = "empty key"
	errKeyContainsNUL   = "invalid key contains NUL character"
	errValueContainsNUL = "invalid value contains NUL character"
)

type ErrInvalidString struct {
	str  string
	text string
}

func newErrInvalidString(s string, text string) error {
	return &ErrInvalidString{s, text}
}

func (e *ErrInvalidString) Arg() string {
	return e.str
}

func (e *ErrInvalidString) Error() string {
	return e.text
}
