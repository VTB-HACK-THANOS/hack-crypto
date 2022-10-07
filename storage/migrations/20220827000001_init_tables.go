package migrations

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var Migration = migrate.NewMigrations(migrate.WithMigrationsDirectory("migrations"))

func init() {

	Migration.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(
			`
			CREATE TABLE users (
			email VARCHAR(50) PRIMARY KEY,
			password TEXT NOT null,
			job_title varchar(255),
            name varchar(255)
			);

			CREATE TABLE roles (
			id integer CHECK (id > 0),
			name varchar(255),
		  PRIMARY KEY (id)
			);

			CREATE TABLE user_roles (
			user_email VARCHAR(50),
			role_id integer,
			PRIMARY KEY (user_email,role_id)
			);

		ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_user_email
					            FOREIGN KEY (user_email)
					            REFERENCES users (email)
					            ON DELETE CASCADE
					            ON UPDATE CASCADE;			           

		ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_role_id
					            FOREIGN KEY (role_id)
					            REFERENCES roles (id)
					            ON DELETE CASCADE
					            ON UPDATE CASCADE;			           

			INSERT INTO public.users (email, "password", job_title, "name") VALUES('tester@test.ru', '$2a$04$oLlN7LTjl8ftPmzq0EXSPugUez7MGPEWJWs0Fnjf1xLQlokmpIs1S', NULL, NULL);

			INSERT INTO public.roles (id, "name") VALUES(1, 'user');
      INSERT INTO public.roles (id, "name") VALUES(2, 'admin');

			INSERT INTO public.user_roles (user_email, role_id) VALUES('tester@test.ru', 2);
     `)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(` 

        `)
		return err
	})
}
