package models

type Gateway struct {
	Request             Request             `json:"request"`
	UpstreamURI         string              `json:"upstream_uri", bson:"upsteam_uri"`
	Response            Response            `json:"response"`
	AuthenticatedEntity AuthenticatedEntity `json:"authenticated_entity"`
	Route               Route               `json:"route"`
	Service             Service             `json:"service"`
	Latencies           Latencies           `json:"latencies"`
	ClientIP            string              `json:"client_ip", bson:"client_ip"`
	// FIX: verify the size of the type
	StartedAt uint32 `json:"started_at"`
}

type Request struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
	URL    string `json:"url"`
	// FIX: verify the size of the type
	Size        uint16        `json:"size"`
	Querystring interface{}   `json:"querystring"`
	Headers     RequestHeader `json:"headers"`
}

type Response struct {
	// FIX: verify the size of the type
	Status uint16 `json:"status"`
	// FIX: verify the size of the type
	Size    uint16         `json:"size"`
	Headers ResponseHeader `json:"headers"`
}

type AuthenticatedEntity struct {
	ConsumerID UUID `json:"consumer_id"`
}

type UUID struct {
	ID string `json:"uuid"`
}

type ID struct {
	ID string `json:"id"`
}

type Timestamp struct {
	// FIX: this type seems strange. ex of input: 1521555129
	CreatedAt int32 `json:"created_at"`
	// FIX: this type seems strange. ex of input: 1521555129
	UpdatedAt int32 `json:"updated_at"`
}

type Route struct {
	Timestamp
	ID
	Hosts        string   `json:"hosts"`
	Methods      []string `json:"methods"`
	Paths        []string `json:"paths"`
	PreserveHost bool     `json:"preserve_host"`
	Protocols    []string `json:"protocols"`
	// FIX: verify the size of the type
	RegexPriority uint64       `json:"regex_priority"`
	Service       RouteService `json:"service"`
	StripPath     bool         `json:"strip_path"`
}

type Service struct {
	Timestamp
	ID
	// FIX: verify the size of the type
	ConnectTimout uint32 `json:"connect_timeout"`
	Host          string `json:"host"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	// FIX: verify the size of the type
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
	// FIX: the type could be not bigger enough
	ReadTimeout uint32 `json:"read_timeout"`
	// FIX: has a maximum number of tries?
	Retries uint16 `json:"retries"`
	// FIX: the type could be not bigger enough
	WriteTimout uint32 `json:"write_timeout"`
}

type Latencies struct {
	// FIX: verify the size of the type
	Proxy uint16 `json:"proxy"`
	// FIX: verify the size of the type
	Gateway uint16 `json:"gateway"`
	// FIX: verify the size of the type
	Request uint16 `json:"request"`
}

type RequestHeader struct {
	Accept    string `json:"accept"`
	Host      string `json:"host"`
	UserAgent string `json:"user-agent"`
}

type ResponseHeader struct {
	ContentLength                 string `json:"content-length"`
	Via                           string `json:"via"`
	Connection                    string `json:"connection"`
	AccessControlAllowCredentials string `json:"access-control-allow-credentials"`
	ContentType                   string `json:"content-type"`
	Server                        string `json:"server"`
	AccessControlAllowOrigin      string `json:"access-control-allow-origin"`
}

type RouteService struct {
	ID
}
