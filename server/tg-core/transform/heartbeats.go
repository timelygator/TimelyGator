package transform

import (
    "log"
    "time"

    "timelygator/server/models"
)

func HeartbeatMerge(lastEvent, heartbeat models.Event, pulsetime float64) *models.Event {
    // Instead of `lastEvent.DataEqual(heartbeat.Data)`,
    // now call `lastEvent.DataEqualJSON(heartbeat.Data)`.
    if lastEvent.DataEqualJSON(heartbeat.Data) {
        // seconds between end of last_event and start of heartbeat
        pulseEnd := lastEvent.Timestamp.Add(lastEvent.Duration).Add(time.Duration(pulsetime) * time.Second)

        withinWindow := lastEvent.Timestamp.Before(heartbeat.Timestamp) && heartbeat.Timestamp.Before(pulseEnd)

        if withinWindow {
            newDuration := heartbeat.Timestamp.Sub(lastEvent.Timestamp) + heartbeat.Duration
            if newDuration < 0 {
                log.Println("Merging heartbeats would result in negative duration, refusing to merge.")
            } else {
                // Extend the lastEvent duration if needed
                if lastEvent.Duration < newDuration {
                    lastEvent.Duration = newDuration
                }
                return &lastEvent
            }
        }
    }
    return nil
}
