# 3-rest-database

This example builds on top of the `2-rest-inmemory` example by refactoring the server to interact with an actual (Postgres) database.

## Requirements

* A running Postgres service is needed
* A `.env` file is used to provide the Postgres database connection parameters
  * Example `.env`:
	  ```
		DB_HOST=localhost
		DB_USER=postgres
		DB_PASSWORD=postgres
		DB_NAME=postgres
		DB_PORT=5432
  * You must use the same environment variable keys shown in the example

## Code Walkthrough

You can find a more in-depth code walk-through in the provided doc [here](code-walkthrough.md).
