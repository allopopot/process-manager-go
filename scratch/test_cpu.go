package main

import (
    "fmt"
    "time"
    "github.com/shirou/gopsutil/cpu"
)

func main() {
    for i := 0; i < 3; i++ {
        overall, _ := cpu.Percent(0, false)
        percore, _ := cpu.Percent(0, true)
        fmt.Println("overall:", overall)
        fmt.Println("percore:", percore)
        time.Sleep(1 * time.Second)
    }
}
