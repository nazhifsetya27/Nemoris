package database

// Init connects to PostgreSQL and runs migrations.
// Single entry point for database boot—keeps main.go free of startup internals.
func Init() {
	Connect()
	Migrate()
}
