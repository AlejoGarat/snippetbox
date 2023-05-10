package formhelpers

import (
	"errors"
	"net/http"

	"github.com/go-playground/form/v4"
)

func DecodePostForm(r *http.Request, dst any, formDecoder *form.Decoder) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}
