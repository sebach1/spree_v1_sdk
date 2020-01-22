package spree

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrNilKey = errors.New("the given api key is nil")

	// ErrNilProductId = errors.New("the given PRODUCT Id is NIL")
	ErrNilProduct = errors.New("the given PRODUCT TITLE is NIL")
	// ErrNilPictures  = errors.New("the given PRODUCT PICTURES are NIL")
	// ErrNilStock     = errors.New("the given PRODUCT STOCK is NIL")

	ErrNilPrice              = errors.New("the given PRODUCT PRICE is NIL")
	ErrNilShippingCategoryId = errors.New("the given PRODUCT SHIPPING CATEGORY ID is NIL")
	ErrNilProductName        = errors.New("the given PRODUCT NAME is NIL")

	ErrStockItemNotFound = errors.New("the given STOCK ITEM was NOT FOUND")

	ErrNilVariant = errors.New("the given VARIANT is NIL")
	// errVariantNotFound = errors.New("the given VARIANT does NOT EXISTS")
	// ErrNilCategory     = errors.New("the given CATEGORY is NIL")
	// errInvalIdListingTypeId = errors.New("the given LISTING TYPE Id is INVALID")

	// errInvalIdCategoryId = errors.New("the given CATEGORY Id is INVALId")
	// ErrNilCategoryId     = errors.New("the given CATEGORY Id is NIL")
	// ErrNilCombinations   = errors.New("the given ATTR COMBINATIONS are NIL")

	// errInvalIdBuyingMode = errors.New("the BUYING MODE is invalid")
	// errInvalIdCondition  = errors.New("the CONDITION is invalid")

	// ErrNilVarStock    = errors.New("the VARIANT wanted to be created has NIL STOCK")
	// ErrNilVarPrice    = errors.New("the VARIANT wanted to be created has NIL PRICE")
	// ErrNilVarPictures = errors.New("the VARIANT wanted to be created has NIL PICTURES")

	// errIncompatibleVar = errors.New("the given VARIANT is INCOMPATIBLE")

	ErrRemoteInconsistency = errors.New("the SERVER had an inconsistency while performing a request (status code != real behaviour)")
)

type Error struct {
	Err    string              `json:"error,omitempty"`
	Errors map[string][]string `json:"errors,omitempty"`
}

var ErrorsSeparator = ";"

func (err *Error) Error() string {
	var strErr string
	for field, validationErrs := range err.Errors {
		strErr += field
		strErr += ":"
		for _, err := range validationErrs {
			strErr += err
			strErr += "  "
		}
		strErr += ErrorsSeparator
	}
	return strErr
}

func errFromReader(stream io.Reader) error {
	body := &Error{}
	err := json.NewDecoder(stream).Decode(body)
	if err != nil {
		return err
	}
	if body.Err == "" {
		body.Err = ErrRemoteInconsistency.Error()
	}
	return body
}
