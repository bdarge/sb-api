
### Database
generate migration script
```console
migrate create -ext sql -dir db/migrations add_account_table
```

#Update go modules
```console
go get -u -d ./...
```