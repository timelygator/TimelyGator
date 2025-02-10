package transform

import (
    "log"
    "time"

    "timelygator/server/models"
)

func HeartbeatMerge(lastEvent, heartbeat models.Event, pulsetime float64) *models.Event {
    if lastEvent.DataEqual(heartbeat.Data) {
        // Seconds between end of last_event and start of heartbeat
        pulseperiodEnd := lastEvent.Timestamp.Add(lastEvent.Duration).Add(time.Duration(pulsetime) * time.Second)
        withinPulsetimeWindow := lastEvent.Timestamp.Before(heartbeat.Timestamp) && heartbeat.Timestamp.Before(pulseperiodEnd)

        if withinPulsetimeWindow {
            // Seconds between end of last_event and start of timestamp
            newDuration := heartbeat.Timestamp.Sub(lastEvent.Timestamp) + heartbeat.Duration
            if lastEvent.Duration < 0 {
                log.Println("Merging heartbeats would result in a negative duration, refusing to merge.")
            } else {
                // Taking the max of durations ensures heartbeats that end before the last event don't shorten it
                if lastEvent.Duration < newDuration {
                    lastEvent.Duration = newDuration
                }
                return &lastEvent
            }
        }
    }
    return nil
}

