package nginx_parser

type Visitor struct {
	IP             string
	IPLocation     string
	RequestTime    string
	RequestMethod  string
	RequestURI     string
	UserAgent      string
	ResponseStatus string
}
