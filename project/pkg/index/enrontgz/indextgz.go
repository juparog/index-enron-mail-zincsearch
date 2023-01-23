package idxenrontgz

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/logger"
	"github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/utils"
)

const (
	MAX_SIZE_PER_FILE = 50_000_000
	MAX_DISK_CAPACITY = 300_000_000
)

var wgRead sync.WaitGroup

type InputParams struct {
	Fields     *string
	Format     *string
	IndexName  *string
	InputPath  *string
	Separator  *string
	Terminator *string
	TokenZinc  *string
	UrlZinc    *string
}

func ProcessFile(inputParams InputParams, done chan<- bool) {
	logger.Info.Println("⚒ Starting processing")
	startTime := time.Now()

	utils.CheckError(utils.CheckFileExist(*inputParams.InputPath))

	file, err := os.Open(*inputParams.InputPath)
	utils.CheckError(err)
	defer file.Close()

	tgz, err := gzip.NewReader(file)
	utils.CheckError(err)
	tr := tar.NewReader(tgz)

	writerChannel := make(chan map[string]string)

	go writeFile(inputParams, writerChannel)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		internalFilepath := header.Name
		switch header.Typeflag {
		case tar.TypeDir: // = directory
			// logger.Debug.Println("Directory:", internalFilepath)
		case tar.TypeReg: // = regular file
			// logger.Debug.Println("Regular file:", internalFilepath)

			buf := new(bytes.Buffer)
			buf.ReadFrom(tr)

			wgRead.Add(1)
			go readBuffer(buf, internalFilepath, inputParams, writerChannel, &wgRead)

		}
	}

	wgRead.Wait()
	close(writerChannel)

	elapsed := time.Since(startTime)
	logger.Info.Printf("🏁 Processing completed in %s\n", elapsed)
	done <- true
}

func readBuffer(buf io.Reader, filepath string, inputParams InputParams, writerChannel chan<- map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(buf)
	scanner.Buffer(utils.GetBufferCapacity())
	// scanner.Split(utils.SplitAt(terminator))

	defaultF := "Body"
	fields := inputParams.Fields
	fieldsMap := makeFieldsMap(inputParams.Separator, fields, &defaultF)
	fieldsMap["Internal-Filepath"] = filepath

	scanner.Scan()
	for {
		line := scanner.Text()
		err := scanner.Err()
		utils.CheckError(err)

		if strings.HasPrefix(line, *inputParams.Terminator) { // finish record, next record
			writerChannel <- utils.CopyMap(fieldsMap)
			clearFieldsMap(&fieldsMap)

			if !scanner.Scan() {
				// close(writerChannel)
				break
			}
			continue
		}

		processLine(&line, inputParams.Separator, &defaultF, &fieldsMap)

		if !scanner.Scan() {
			writerChannel <- utils.CopyMap(fieldsMap)
			// close(writerChannel)
			break
		}
	}
}

func makeFieldsMap(separator *string, fields *string, defaultF *string) map[string]string {
	fieldsResult := make(map[string]string)
	fieldsResult[*defaultF] = ""

	fieldsList := strings.Split(*fields, *separator)
	for _, field := range fieldsList {
		fieldsResult[field] = ""
	}

	return fieldsResult
}

func clearFieldsMap(fieldsMap *map[string]string) {
	for fieldName := range *fieldsMap {
		(*fieldsMap)[fieldName] = ""
	}
}

func processLine(line *string, separator *string, defaultF *string, fieldsMap *map[string]string) {
	dataList := strings.Split(*line, *separator)
	for fieldName := range *fieldsMap {
		if fieldName == dataList[0] {
			// concatenate all the information that follows in the line to the field,
			// it is done with a for because the line can contain the separator more than once
			for i := 1; i < len(dataList); i++ {
				(*fieldsMap)[fieldName] = dataList[i]
			}
			return
		}
	}

	(*fieldsMap)[*defaultF] = (*fieldsMap)[*defaultF] + *line
	return
}

func writeFile(inputParams InputParams, writerChannel <-chan map[string]string) {
	logger.Info.Println("Writing files...")

	writeBasePath := utils.GetBasePath(*inputParams.InputPath)
	fileName := utils.GetNewFileName(*inputParams.InputPath)
	fileExt := utils.GetExtension(*inputParams.Format)
	partCounter := 1
	filesUploading := 0
	tmpFilepath := filepath.Join(writeBasePath, fmt.Sprint(fileName, "-Part-", partCounter, fileExt))
	writeStringFunc := createStringWriter(tmpFilepath)
	jsonFunc := utils.GetJSONFunc(*inputParams.Format)

	if *inputParams.IndexName == "" {
		inputParams.IndexName = &fileName
	}
	beforeIdxStr, middleIdxStr, afterIdxStr := utils.MakeIndexString(inputParams.IndexName, inputParams.Format)

	record, more := <-writerChannel
	writeStringFunc(beforeIdxStr, false)
	for more {

		jsonData := jsonFunc(record)
		currentSize := writeStringFunc(jsonData, false)

		record, more = <-writerChannel
		if more {
			if currentSize > MAX_SIZE_PER_FILE {
				writeStringFunc(afterIdxStr, true)
				logger.Info.Printf("write complete: %s", tmpFilepath)

				go uploadFile(tmpFilepath, inputParams, &filesUploading)

				for {
					if (filesUploading * MAX_SIZE_PER_FILE) <= MAX_DISK_CAPACITY {
						break
					}
					time.Sleep(1 * time.Second)
				}

				partCounter += 1
				tmpFilepath = filepath.Join(writeBasePath, fmt.Sprint(fileName, "-", partCounter, fileExt))
				writeStringFunc = createStringWriter(tmpFilepath)
				writeStringFunc(beforeIdxStr, false)
			} else {
				writeStringFunc(middleIdxStr, false)
			}
		} else {
			break
		}
	}

	writeStringFunc(afterIdxStr, true)
	logger.Info.Printf("write complete: %s", tmpFilepath)

	go uploadFile(tmpFilepath, inputParams, &filesUploading)

	logger.Info.Println("Writing completed.")
}

func createStringWriter(filePath string) func(string, bool) int64 {
	file, err := os.Create(filePath)
	utils.CheckError(err)
	logger.Debug.Printf("File created: %s\n", filePath)

	return func(data string, close bool) int64 {
		_, err := file.WriteString(data)
		utils.CheckError(err)

		if close {
			file.Close()
			logger.Debug.Printf("File closed: %s\n", filePath)
			return 0
		}

		stat, err := file.Stat()
		utils.CheckError(err)

		return stat.Size()
	}
}

func uploadFile(filepath string, inputParams InputParams, filesUploading *int) {
	*filesUploading += 1
	utils.CheckError(utils.CheckFileExist(filepath))

	file, err := os.Open(filepath)
	utils.CheckError(err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	logger.Info.Printf("Uploading: %s\n", filepath)
	targetUrl := *inputParams.UrlZinc + utils.GetUploadEndpoint(*inputParams.Format)
	req, err := http.NewRequest("POST", targetUrl, buf)
	utils.CheckError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+*inputParams.TokenZinc)

	client := &http.Client{}
	res, err := client.Do(req)
	utils.CheckError(err)

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		utils.CheckError(fmt.Errorf("Error uploading file %s, api response with status code %d, response: %s\n", filepath, res.StatusCode, string(body)))
	}

	logger.Info.Printf("Upload complete: %s, response: %s", filepath, string(body))

	*filesUploading -= 1

	file.Close()
	utils.CheckError(os.Remove(filepath))
}
