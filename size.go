package nzb

import (
    "fmt"
)

// Size .
type Size uint64

func (s Size) String() string {
    sizeType := "Byte(s)"
    scnt := 0
    size := float64(s)
    for size > 1024 {
        size /= float64(1024)
        scnt++
    }

    switch scnt {
    case 1:
        sizeType = "KB"
    case 2:
        sizeType = "MB"
    case 3:
        sizeType = "GB"
    case 4:
        sizeType = "TB"
    }

    return fmt.Sprintf("%.2f %s", size, sizeType)
}

// EOF
