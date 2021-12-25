package database

type Dsn struct {
	host     string
	port     string
	dbname   string
	user     string
	password string
}

func NewDsn(host string, port string, database string, user string, password string) *Dsn {
	d := Dsn{
		host:     host,
		port:     port,
		dbname:   database,
		user:     user,
		password: password,
	}
	return &d

}
func (d *Dsn) Dsn() string {
	return "postgres://" + d.user + ":" + d.password + "@" + d.host + ":" + d.port + "/" + d.dbname + "?sslmode=disable&TimeZone=Asia/Tokyo"
}
