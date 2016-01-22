package main

import (
  "fmt"
  "os"
)


func main() {

  init()

}

func init() {
  // check to see if config file exists
  if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
    fmt.Println("file does not exist")
  }
}
