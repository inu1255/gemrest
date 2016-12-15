package gemrest

type database struct {
	Driver     string
	DataSource string
}

type config struct {
	Database database
	Dev      bool
	DocBind  string // doc api url
}
