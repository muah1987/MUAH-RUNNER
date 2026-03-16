package questions

import (
"fmt"
"time"
)

type CountdownTimer struct {
Duration time.Duration
}

func NewCountdownTimer(seconds int) *CountdownTimer {
return &CountdownTimer{Duration: time.Duration(seconds) * time.Second}
}

func (t *CountdownTimer) Display(onTick func(remaining int)) {
remaining := int(t.Duration.Seconds())
for remaining > 0 {
onTick(remaining)
time.Sleep(1 * time.Second)
remaining--
}
fmt.Print("\r")
}
