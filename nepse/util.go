package nepse

import "time"

func minInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func minDuration(a, b time.Duration) time.Duration {
    if a < b {
        return a
    }
    return b
}

