package goose

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"
)

// Status prints the status of all migrations.
func Status(db *sql.DB, dir string) error {
	// collect all migrations
	migrations, err := CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}

	// must ensure that the version table exists if we're running on a pristine DB
	if _, err := EnsureDBVersion(db); err != nil {
		return err
	}

	printMigrationStatus(db, migrations)

	return nil
}

func printMigrationStatus(db *sql.DB, migrations Migrations) {
	records := make(map[int64]MigrationRecord)
	q := fmt.Sprintf("SELECT version_id, tstamp, is_applied FROM %s ORDER BY tstamp ASC", TableName())
	rows, err := db.Query(q)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	for rows.Next() {
		var row MigrationRecord
		rows.Scan(&row.VersionID, &row.TStamp, &row.IsApplied)
		records[row.VersionID] = row
	}

	log.Println("    Applied At                  Migration")
	log.Println("    =======================================")
	for _, migration := range migrations {
		var appliedAt string
		if row, ok := records[migration.Version]; ok && row.IsApplied {
			appliedAt = row.TStamp.Format(time.ANSIC)
		} else {
			appliedAt = "Pending"
		}
		log.Printf("    %-24s -- %v\n", appliedAt, filepath.Base(migration.Source))
	}
}
