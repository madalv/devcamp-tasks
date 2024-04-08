All SQL statements can be found in `statements.sql`.

To create the database, simply run `make run`, then `make migup`, which will run the migrations on the MySQL database.

`db/migrations` currently holds the first migration that defines the initial schema.

The code that seeds the database is in `util/seeder.go`.