package background

import (
	"database/sql"
	"log"
	"time"
)

func DeleteExpiredSessions(db *sql.DB) {
	_, err := db.Exec("DELETE FROM sessions WHERE expire_time < DATETIME('now')")
	if err != nil {
		log.Fatal("Error deleting expired sessions:", err)
	}
}

func WorkerScanBD(db *sql.DB) {
	ticker := time.NewTicker(10 * time.Minute) // adjust the interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			DeleteExpiredSessions(db)
		}
	}
}
