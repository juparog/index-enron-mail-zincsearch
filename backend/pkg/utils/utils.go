package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	idxtypes "github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/index/types"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func exitGracefully(err error) {
	log.Fatal(err)
	os.Exit(1)
}

func CheckError(e error) {
	if e != nil {
		exitGracefully(e)
	}
}

func IfThenElse[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}

func CheckIfValidFile(filename string) error {
	if fileExtension := filepath.Ext(filename); fileExtension != ".tgz" {
		return fmt.Errorf("File %s is not TGZ", filename)
	}

	return CheckFileExist(filename)
}

func CheckFileExist(filename string) error {
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("File %s does not exist", filename)
	}

	return nil
}

func GetBufferCapacity() ([]byte, int) {
	const maxCapacity = 512 * 1024
	buf := make([]byte, maxCapacity)

	return buf, maxCapacity
}

func CopyMap[K, V comparable](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func GetBasePath(strFilepath string) string {
	return filepath.Dir(strFilepath)
}

func GetNewFileName(strFilepath string) string {
	return strings.TrimSuffix(filepath.Base(strFilepath), filepath.Ext(strFilepath))
}

func GetExtension(format string) string {
	fileExt := ".json"
	switch format {
	case idxtypes.Z_BULK_FORMAT:
		fileExt = ".ndjson"
	}
	return fileExt
}

func MakeIndexString(indexName *string, format *string) (string, string, string) {
	beforeIndexStr := "[\n\t"
	middleIndexStr := ",\n\t"
	afterIndexStr := "\n]"
	switch *format {
	case idxtypes.Z_BULK_FORMAT:
		beforeIndexStr = `{ "index": { "_index": "` + *indexName + `" } }
`
		middleIndexStr = `
{ "index": { "_index": "` + *indexName + `" } }
`
		afterIndexStr = ""
	case idxtypes.Z_BULKV2_FORMAT:
		beforeIndexStr = `{
	"index": "` + *indexName + `",
	"records": [
		`
		middleIndexStr = `,
		`
		afterIndexStr = `
	]
}`
	}

	return beforeIndexStr, middleIndexStr, afterIndexStr
}

func GetJSONFunc(format string) func(map[string]string) string {
	switch format {
	case idxtypes.Z_BULK_FORMAT:
		return func(record map[string]string) string {
			jsonData, _ := json.Marshal(record)
			return string(jsonData)
		}
	case idxtypes.Z_BULKV2_FORMAT:
		return func(record map[string]string) string {
			jsonData, _ := json.Marshal(record)
			return string(jsonData)
		}
	default:
		return func(record map[string]string) string {
			jsonData, _ := json.MarshalIndent(record, "\t", "\t")
			return string(jsonData)
		}
	}
}

func GetUploadEndpoint(format string) string {
	endpoint := "_bulk"
	switch format {
	case idxtypes.Z_BULKV2_FORMAT:
		endpoint = "_bulkv2"
	}
	return endpoint
}
