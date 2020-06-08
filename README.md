# example-go-rest-api
Demonstrates different ways to build a golang RESTful API from basic net/http to using an ORM.

## Steps

We'll walk through and iterate on building a RESTful API service written in Go like so:
1. `1-hello-world`: Basic web server that outputs "Hello world!" when navigating to the root url.
2. `2-rest-inmemory`: Web server that implements CRUD operations against an in-memory "database".
3. `3-rest-database`: Builds off step 2 by replacing in-memory "database" with a real SQL-based database.
