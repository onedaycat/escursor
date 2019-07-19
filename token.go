package escursor

import (
	"encoding/base64"
	"errors"
)

type SortValueHandler func(index int) []interface{}
type SliceLastItemHandler func(index int)

func CreateNextToken(size, length int, sortValueHandler SortValueHandler, sliceLastItemHandler SliceLastItemHandler) (string, error) {
	var lastSortValue []interface{}

	if size == 0 || length == 0 {
		return "", nil
	}

	if size != 0 && length < size {
		sliceLastItemHandler(length - 1)
		lastSortValue = sortValueHandler(length - 1)
	} else if size != 0 && length > size {
		sliceLastItemHandler(size - 1)
		lastSortValue = sortValueHandler(size - 1)
	}

	return createNextToken(size, length, lastSortValue)
}

func createNextToken(size, length int, sortedValues []interface{}) (string, error) {
	if length <= size || size == 0 {
		return "", nil
	}

	cursorFields := make(Cursorfields, 0, len(sortedValues))
	for i := 0; i < len(sortedValues); i++ {
		cursorFields = append(cursorFields, sortedValues[i])
	}

	cfByte, err := cursorFields.MarshalMsg(nil)
	if err != nil {
		return "", errors.New("unable create next token")
	}

	return base64.URLEncoding.EncodeToString(cfByte), nil
}
