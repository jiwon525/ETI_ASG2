module hello

go 1.17

require (
	github.com/go-sql-driver/mysql v1.6.0 // direct
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0 // direct
)

require (
	github.com/felixge/httpsnoop v1.0.1 // indirect
	gorm.io/gorm v1.22.5
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
)

//cd GO-WORK\src\ETI_ASG2\DB\main.go
//go mod init "api url"
//e.g. MarksWallet.com/api/v1/Token
