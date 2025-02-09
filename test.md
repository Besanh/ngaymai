Run in cmd
`go install github.com/golang/mock/mockgen@latest`

`go get github.com/golang/mock/gomock`

`mockgen -source=common/cache/redis.go -destination=mock/mock_redis.go -package=mock`

`mockgen -source=common/sqlclient/sql.go -destination=mock/mock_db.go -package=mock`

`go generate ./...`