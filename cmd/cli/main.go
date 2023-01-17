package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	idxenrontgz "github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/index/enrontgz"
	idxtypes "github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/index/types"
	"github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/utils"
)

var formatList = []string{idxtypes.Z_BULK_FORMAT, idxtypes.Z_BULKV2_FORMAT, idxtypes.JSON_FORMAT}

func main() {

	if len(os.Args) < 2 {
		utils.CheckError(fmt.Errorf("argument <filepath> is required."))
	}

	flag.Usage = func() {
		fmt.Printf("Usage: [options] %s <filePath>\nOptions:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	fields := flag.String("fields", "", "data fields. (default: all content in the 'Body' field).")
	format := flag.String("format", idxtypes.Z_BULK_FORMAT, fmt.Sprintf("output format file, valid: %v.", formatList))
	indexName := flag.String("index-name", "", "index name in ZinSearch.")
	separator := flag.String("separator", ":", "character separating each field in a record.")
	terminator := flag.String("terminator", "\r\n", "character separating each row in data.")
	token := flag.String("token", "", fmt.Sprintf("authorization basic token for api ZincSearch."))
	urlZinc := flag.String("url-zinc", "http://localhost:4080/api/", "ZincSearch url to index.")
	verbosity := flag.Bool("verbosity", false, "verbosity log.")

	flag.Parse()

	if !(utils.Contains(formatList, *format)) {
		utils.CheckError(fmt.Errorf("format not supported, valid: %v.", formatList))
	}

	if len(*token) == 0 {
		utils.CheckError(fmt.Errorf("token cannot be empty"))
	}

	inputPath := flag.Arg(0)

	tgzParams := idxenrontgz.InputParams{
		Fields:     fields,
		Format:     format,
		IndexName:  indexName,
		InputPath:  &inputPath,
		Separator:  separator,
		Terminator: terminator,
		TokenZinc:  token,
		UrlZinc:    urlZinc,
		Verbosity:  verbosity,
	}

	done := make(chan bool)

	go idxenrontgz.ProcessFile(tgzParams, done)

	<-done

}
