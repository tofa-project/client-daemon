package db

// Creates default apps table
func makeAppsTable() {
	Instance.Exec(`
		CREATE TABLE apps (
			"id" TEXT UNIQUE,
			"data" BLOB			
		);
	`)
}

// Creates default logs table
func makeLogsTable() {
	Instance.Exec(`
		CREATE TABLE logs (
			"id" TEXT,
			"data" BLOB			
		);
	`)
}
