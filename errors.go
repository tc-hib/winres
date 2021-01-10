package winres

const (
	errZeroID          = "ordinal identifier must not be zero"
	errEmptyName       = "string identifier must not be empty"
	errNameContainsNUL = "string identifier must not contain NUL char"

	errUnknownArch = "unknown architecture"

	errNotICO                 = "not a valid ICO file"
	errImageLengthTooBig      = "image size found in ICONDIRENTRY is too big (above 10 MB)"
	errTooManyIconSizes       = "too many sizes"
	errGroupNotFound          = "group does not exist"
	errInvalidGroup           = "invalid group"
	errIconMissing            = "icon missing from group"
	errCursorMissing          = "cursor missing from group"
	errInvalidImageDimensions = "invalid image dimensions"
	errImageTooBig            = "image size too big, must fit in 256x256"
	errNotCUR                 = "not a valid CUR file"
	errUnknownImageFormat     = "unknown image format"
)
