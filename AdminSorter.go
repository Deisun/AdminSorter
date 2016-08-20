// Author:  Rob Douma
// Description: This program will monitor a folder for worksheet files and automatically sort them into Eworksheets folders

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

/*var destinationPath string = "c:\\temp1\\Eworkz"
var sourcePath string = "c:\\temp1\\ToSort"*/

var sourcePath, err = os.Getwd()
var destinationPath string = "\\\\usatfs01\\Eworksheets"

func main() {

	Run()

}

func Run() {
	// get a list of all the files in the dir
	fileList := GetFiles()
	// create a list of only filenames that are integers above 1,000,000.
	// this is done because there are some funky foldernames for worksheets under that value
	worksheetList := ExtractValidWorksheets(fileList)

	GenerateDirs(worksheetList)
	MoveWorksheets(worksheetList)
}

func GetFiles() []string {
	// get a list of all files in the Dir
	files, err := ioutil.ReadDir(sourcePath)

	if err != nil {
		log.Fatal(err)
	}

	// create a dynamic slice of filenames
	fileList := make([]string, 0)

	for _, file := range files {
		name := file.Name()

		fileList = append(fileList, name)
	}

	return fileList
}

func ExtractValidWorksheets(fileList []string) []string {
	worksheetList := make([]string, 0)

	for _, file := range fileList {
		// remove the last 4 characters to remove .ext
		fileTrimmed := strings.TrimSuffix(file, ".pdf")

		// check to see if integer in filename
		if _, err := strconv.Atoi(fileTrimmed); err == nil {
			i, err := strconv.Atoi(fileTrimmed)

			if err != nil {
				log.Fatal(err)
			}

			if i >= 700001 {
				worksheetList = append(worksheetList, file)
			}
		}
	}

	return worksheetList
}

// creates the directory folders if they don't exist
func GenerateDirs(worksheetList []string) {

	for _, worksheet := range worksheetList {
		// remove the last 4 characters to remove .ext
		worksheetInt := ConvertWorksheetStringToInt(worksheet)

		worksheetFolderName := GetWorksheetsFolderName(worksheetInt)
		worksheetFullPath := destinationPath + "\\" + worksheetFolderName

		if _, err := os.Stat(worksheetFullPath); err == nil {
			// do nothing, path already exists
		} else {

			err := os.Mkdir(worksheetFullPath, 0777)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
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
func GetWorksheetsFolderName(worksheetNum int) string {

	if worksheetNum%500 == 0 {
		worksheetNum -= 1
	}

	firstNum := worksheetNum/500*500 + 1
	secondNum := firstNum + 499

	folderName := strconv.Itoa(firstNum) + "-" + strconv.Itoa(secondNum)

	return folderName
}

func ConvertWorksheetStringToInt(worksheet string) int {
	worksheetTrimmed := strings.TrimSuffix(worksheet, ".pdf")

	i, err := strconv.Atoi(worksheetTrimmed)

	if err != nil {
		log.Fatal(err)
	}
	return i
}

func MoveWorksheets(worksheetList []string) {

	var wg sync.WaitGroup

	for _, worksheet := range worksheetList {

		wg.Add(1)

		worksheetInt := ConvertWorksheetStringToInt(worksheet)

		src := sourcePath + "\\" + worksheet
		dst := destinationPath + "\\" + GetWorksheetsFolderName(worksheetInt) + "\\" + worksheet

		go func() {
			err := CopyFile(src, dst)

			if err != nil {
				log.Fatal(err)
			} else {
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
