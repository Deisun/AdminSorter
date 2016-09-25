// Author:  Rob Douma
// Description: This program moves all worksheet files into Eworksheets shared folders using the existing folder naming scheme.
//				If the folder does not exist, it will automatically create a folder according to the worksheets name

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Worksheet struct {
	destinationFolder string
	name              string
	value             int
}

// var destinationPath string = "c:\\temp1\\Eworkz"
var destinationPath string = "\\\\usatfs01\\Eworksheets"
var sourcePath, err = os.Getwd()

var worksheets = []*Worksheet{}

func main() {
	GetWorksheets()
	GenerateDirs()
	MoveWorksheets()
	// PrintWorksheets()
}

func GetWorksheets() {

	files, err := ioutil.ReadDir(sourcePath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if IsPDF(file) && IsNumber(file) {

			AppendToWorksheets(file)
		}
	}
}

func IsPDF(file os.FileInfo) bool {
	buf, _ := ioutil.ReadFile(file.Name())

	if len(buf) > 3 && // these correspond to PDF "magic numbers" header bytes for %PDF
		buf[0] == 0x25 && // hex: %
		buf[1] == 0x50 && // hex: P
		buf[2] == 0x44 && // hex: D
		buf[3] == 0x46 { // hex: F
		return true
	}

	return false
}

func IsNumber(file os.FileInfo) bool {
	filenameNoExtension := strings.TrimSuffix(file.Name(), ".pdf")

	_, err := strconv.Atoi(filenameNoExtension)

	if err != nil {
		return false
	}

	return true
}

func AppendToWorksheets(file os.FileInfo) {
	worksheet := new(Worksheet)

	worksheet.name = file.Name()
	worksheet.value = GetFileIntValue(file)
	worksheet.destinationFolder = GetDestinationFolderName(worksheet.value)

	worksheets = append(worksheets, worksheet)
}

func GetFileIntValue(file os.FileInfo) int {
	filenameNoExtension := strings.TrimSuffix(file.Name(), ".pdf")

	value, err := strconv.Atoi(filenameNoExtension)

	if err != nil {
		log.Fatal(err)
	}

	return value
}

// GetEworksheetsFolder returns a folder name based off the filename of a worksheet.
// The syntax of an Eworksheets folders is:  "####-####", where the first number begins at 1
// and the second number is 499 units apart from the first.
// Example:  "1-500", "501-1000", "1001-1500", "1501-2000", etc.

// To achieve the correct folder name, the worksheet filename is first converted to an integer.
// It is then divided by 500.  We use an integer so that the decimals are ignored after the division.
// Ignoring the decimals allows us to then multiply by 500 and add 1 to get the first number of the folder
// Example:  1035529 / 500 = 2071.058, but since we're diving integers we only get 2071
// then multiplying by 500 and adding 1 gives us 1035501
func GetDestinationFolderName(worksheetNum int) string {

	if worksheetNum%500 == 0 {
		worksheetNum -= 1
	}

	firstNum := worksheetNum/500*500 + 1
	secondNum := firstNum + 499

	folderName := strconv.Itoa(firstNum) + "-" + strconv.Itoa(secondNum)

	return folderName
}

func GenerateDirs() {

	for _, worksheet := range worksheets {

		worksheetFullPath := destinationPath + "\\" + worksheet.destinationFolder

		if _, err := os.Stat(worksheetFullPath); err == nil {
			// do nothing, path already exists
		} else {

			err := os.Mkdir(worksheetFullPath, 0777)

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Folder does not exist.  Creating folder: " + "\"" + worksheetFullPath + "\"")
			}
		}
	}
}

// func PrintWorksheets() {
// 	for _, worksheet := range worksheets {
// 		fmt.Print("dst: ", worksheet.destinationFolder, "\tname: ", worksheet.name, "\tvalue: ", worksheet.value, "\n")
// 	}
// }

func MoveWorksheets() {

	var wg sync.WaitGroup

	for _, worksheet := range worksheets {

		wg.Add(1)

		src := sourcePath + "\\" + worksheet.name
		dst := destinationPath + "\\" + worksheet.destinationFolder + "\\" + worksheet.name

		go func() {
			err := CopyFile(src, dst)

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Moving: " + src + " to: " + "\"" + dst + "\"")

				rerr := os.Remove(src)

				if rerr != nil {
					log.Fatal(err)
				}
			}

			wg.Done()
		}()

	}

	wg.Wait()
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {

	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
