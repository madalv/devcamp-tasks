All SQL statements can be found in `statements.sql`.

To create the database, simply run `make run`, then `make migup`, which will run the migrations on the MySQL database.

`db/migrations` currently holds the first migration that defines the initial schema. `sqlc` was used to convert the queries from `db/queries` to structs because no one has time to write that stuff manually.

The code that seeds the database is in `db/sqlc/seeder`.