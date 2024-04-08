### Task 1

All SQL statements can be found in `statements.sql`.

To create the database, simply run `make run`, then `make migup`, which will run the migrations on the MySQL database.

`db/migrations` currently holds the first migration that defines the initial schema.

The code that seeds the database is in `util/seeder.go`.

### Task 2

Set up a simple router (`api/handler.go`) and a very simple cache (`repository/local_cache.go`). 
There is a benchmark for the `getCampaignsForSource` handler in `api/handler_test.go`.

Run `make bench` to run the test.

Result for the benchmark with the cache:
```
  657496              1649 ns/op
PASS
ok      adt/api 2.129s
```
Results for the benchmark without the cache:
```
    1858            647353 ns/op
PASS
ok      adt/api 1.275s
```