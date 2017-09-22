 package main

import (
  "os"
  "github.com/mbesancon/tomate"
)

 func main() {
   tmt := tomate.New(20, 5, 7, 3)
   tmt.Loop(os.Stdout)
 }
