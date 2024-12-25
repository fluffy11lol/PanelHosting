package modelsControlPanel

import (
	"os"
)

// Чтение переменных окружения
type Envs struct {
	Host_psql      string
	Port_psql      string
	User_psql      string
	Password_psql  string
	Dbname_psql    string
	Host_mysql     string
	Port_mysql     string
	User_mysql     string
	Password_mysql string
	Dbname_mysql   string
	Grpc_port      string
	Rest_port      string
}

var EnvsVars Envs

func LoadEnvs() {
	EnvsVars = Envs{
		Host_psql:      os.Getenv("POSTGRES_HOST"),
		Port_psql:      os.Getenv("POSTGRES_PORT"),
		User_psql:      os.Getenv("POSTGRES_USER"),
		Password_psql:  os.Getenv("POSTGRES_PASSWORD"),
		Dbname_psql:    os.Getenv("POSTGRES_DB"),
		Host_mysql:     os.Getenv("MYSQL_HOST"),
		Port_mysql:     os.Getenv("MYSQL_PORT"),
		User_mysql:     os.Getenv("MYSQL_USER"),
		Password_mysql: os.Getenv("MYSQL_PASSWORD"),
		Dbname_mysql:   os.Getenv("MYSQL_DBNAME"),
		Grpc_port:      os.Getenv("GRPC_PORT"),
		Rest_port:      os.Getenv("REST_PORT"),
	}
}
