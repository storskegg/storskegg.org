package main

import (
    "github.com/storskegg/storskegg.org/application/beardsrv"
    "log"
)

func main() {
    if err := beardsrv.Run(); err != nil {
        log.Fatal(err)
    }
}
