package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {

	Migration.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(
			`
			CREATE TABLE nfts (
			id uuid primary key,
			user_email varchar(50),
			type text
			);

		ALTER TABLE nfts ADD CONSTRAINT fk_nfts_user_email
					            FOREIGN KEY (user_email)
					            REFERENCES users(email)
					            ON DELETE CASCADE
					            ON UPDATE CASCADE;			           

		CREATE TABLE tasks (
			id uuid primary key,
			name varchar(255),
			description text,
			user_email varchar(50),
			type text
			);

		ALTER TABLE tasks ADD CONSTRAINT fk_tasks_user_email
					            FOREIGN KEY (user_email)
					            REFERENCES users(email)
					            ON DELETE CASCADE
					            ON UPDATE CASCADE;			           
     `)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(` 

        `)
		return err
	})
}
