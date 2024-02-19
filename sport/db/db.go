package db

import (
	"fmt"
	"time"

	"syreclabs.com/go/faker"
)

func (r *eventsRepo) seed() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, name TEXT, type TEXT, location TEXT, visible INTEGER, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		sportType := faker.RandomChoice([]string{"Running", "Walking", "Jumping", "Throwing"})
		team := faker.Team()

		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO events(id, name, type, location, visible, advertised_start_time) VALUES (?,?,?,?,?,?)`)
		if err == nil {

			_, err = statement.Exec(
				i,
				fmt.Sprintf("%s %s", team.Name(), sportType),
				sportType,
				team.State(),
				faker.Number().Between(0, 1),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}
