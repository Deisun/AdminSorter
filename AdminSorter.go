// Author:  Rob Douma

package main

import (
  "fmt"
  "os"
)


func main() {

  initialize()

}

func initialize() {
  // check to see if config file exists
  if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
    fmt.Println("file does not exist")
  }
}
