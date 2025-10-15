# sb api

An api based on proto files from [sb-proto](git@github.com:bdarge/sb-proto.git) to do CRUD operations.

Use https://github.com/golang-migrate/migrate/tree/master/cmd/migrate to generate UP/DOWN migration files

```
$ migrations % migrate create -ext sql add_lang
/Users/bdarge/projects/sb/api/db/migrations/20251015052434_add_lang.up.sql
/Users/bdarge/projects/sb/api/db/migrations/20251015052434_add_lang.down.sql
```

