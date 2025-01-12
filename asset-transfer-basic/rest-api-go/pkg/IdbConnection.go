package pkg

type IDBConnection interface {
	Connect() any
}
type DBConnection struct {
	Db IDBConnection
}

func (con DBConnection) DBConnect() any {
	return con.Db.Connect()
}
