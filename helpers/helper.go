package helpers

import (
	"fmt"
	"os"
)

//ConvertToByteArray converst file to byte array
func ConvertToByteArray(filepath string) []byte {
	byteArr, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}

	return byteArr
}
