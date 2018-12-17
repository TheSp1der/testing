package main

import (
        "github.com/nsf/termbox-go"
        "time"
)


func main() {
        // initalize termbox
        err := termbox.Init()
        if err != nil {
                panic(err)
        }
        defer termbox.Close()

        // clear buffer
        termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

        // write message to buffer
        x := 0

        for _, letter := range "Hello World" {
                termbox.SetCell(x, 0, letter, termbox.ColorMagenta, termbox.ColorDefault)
                x++
        }

        // output buffer
        termbox.Flush()

        time.Sleep(time.Duration(time.Second * 5))
}
