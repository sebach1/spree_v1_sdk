package spree

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	errNilKey = errors.New("the given api key is nil")

	errNilProductId = errors.New("the given PRODUCT Id is NIL")
	errNilProduct   = errors.New("the given PRODUCT TITLE is NIL")
	errNilPictures  = errors.New("the given PRODUCT PICTURES are NIL")
	errNilStock     = errors.New("the given PRODUCT STOCK is NIL")

	errNilPrice              = errors.New("the given PRODUCT PRICE is NIL")
	errNilShippingCategoryId = errors.New("the given PRODUCT SHIPPING CATEGORY ID is NIL")
	errNilProductName        = errors.New("the given PRODUCT NAME is NIL")

	errStockItemNotFound = errors.New("the given STOCK ITEM was NOT FOUND")

	errNilVariant      = errors.New("the given VARIANT is NIL")
	errVariantNotFound = errors.New("the given VARIANT does NOT EXISTS")
	// errNilCategory     = errors.New("the given CATEGORY is NIL")
	errInvalIdListingTypeId = errors.New("the given LISTING TYPE Id is INVALID")

	errInvalIdCategoryId = errors.New("the given CATEGORY Id is INVALId")
	errNilCategoryId     = errors.New("the given CATEGORY Id is NIL")
	errNilCombinations   = errors.New("the given ATTR COMBINATIONS are NIL")

	errInvalIdBuyingMode = errors.New("the BUYING MODE is invalid")
	errInvalIdCondition  = errors.New("the CONDITION is invalid")

	errNilVarStock    = errors.New("the VARIANT wanted to be created has NIL STOCK")
	errNilVarPrice    = errors.New("the VARIANT wanted to be created has NIL PRICE")
	errNilVarPictures = errors.New("the VARIANT wanted to be created has NIL PICTURES")

	errIncompatibleVar = errors.New("the given VARIANT is INCOMPATIBLE")

	errRemoteInconsistency = errors.New("the SERVER had an inconsistency while performing a request (status code != real behaviour)")
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
		body.Err = errRemoteInconsistency.Error()
	}
	return body
}
