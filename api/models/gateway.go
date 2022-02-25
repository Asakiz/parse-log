package models

type Gateway struct {
	Request             Request             `json:"request" bson:"request"`
	UpstreamURI         string              `json:"upstream_uri" bson:"upstream_uri"`
	Response            Response            `json:"response" bson:"response"`
	AuthenticatedEntity AuthenticatedEntity `json:"authenticated_entity" bson:"authenticated_entity"`
	Route               Route               `json:"route" bson:"route"`
	Service             Service             `json:"service" bson:"service"`
	Latencies           Latencies           `json:"latencies" bson:"latencies"`
	ClientIP            string              `json:"client_ip" bson:"client_ip"`
	StartedAt           uint32              `json:"started_at"`
}

type Request struct {
	Method      string        `json:"method" bson:"method"`
	URI         string        `json:"uri" bson:"uri"`
	URL         string        `json:"url" bson:"url"`
	Size        uint16        `json:"size" bson:"size"`
	Querystring interface{}   `json:"querystring" bson:"querystring"`
	Headers     RequestHeader `json:"headers" bson:"headers"`
}

type RequestHeader struct {
	Accept    string `json:"accept" bson:"accept"`
	Host      string `json:"host" bson:"host"`
	UserAgent string `json:"user-agent" bson:"user_agent"`
}

type Response struct {
	Status  uint16         `json:"status" bson:"status"`
	Size    uint16         `json:"size" bson:"size"`
	Headers ResponseHeader `json:"headers" bson:"headers"`
}

type ResponseHeader struct {
	ContentLength                 string `json:"content-length" bson:"content_length"`
	Via                           string `json:"via" bson:"via"`
	Connection                    string `json:"connection" bson:"connection"`
	AccessControlAllowCredentials string `json:"access-control-allow-credentials" bson:"access_control_allow_credentials"`
	ContentType                   string `json:"content-type" bson:"content_type"`
	Server                        string `json:"server" bson:"server"`
	AccessControlAllowOrigin      string `json:"access-control-allow-origin" bson:"access_control_allow_origin"`
}

type AuthenticatedEntity struct {
	ConsumerID UUID `json:"consumer_id" bson:"consumer_id"`
}

type UUID struct {
	ID string `json:"uuid" bson:"uuid"`
}

type Route struct {
	CreatedAt     int32        `json:"created_at" bson:"created_at"`
	Hosts         string       `json:"hosts" bson:"hosts"`
	ID            string       `json:"id" bson:"id"`
	Methods       []string     `json:"methods" bson:"methods"`
	Paths         []string     `json:"paths" bson:"paths"`
	PreserveHost  bool         `json:"preserve_host" bson:"preserve_host"`
	Protocols     []string     `json:"protocols" bson:"protocols"`
	RegexPriority uint64       `json:"regex_priority" bson:"regex_priority"`
	Service       RouteService `json:"service" bson:"service"`
	StripPath     bool         `json:"strip_path" bson:"strip_path"`
	UpdatedAt     int32        `json:"updated_at" bson:"updated_at"`
}

type RouteService struct {
	ID string `json:"id" bson:"id"`
}

type Service struct {
	ConnectTimout uint32 `json:"connect_timeout" bson:"connect_timout"`
	CreatedAt     int32  `json:"created_at" bson:"created_at"`
	Host          string `json:"host" bson:"host"`
	ID            string `json:"id" bson:"id"`
	Name          string `json:"name" bson:"name"`
	Path          string `json:"path" bson:"path"`
	Port          uint16 `json:"port" bson:"port"`
	Protocol      string `json:"protocol" bson:"protocol"`
	ReadTimeout   uint32 `json:"read_timeout" bson:"read_timeout"`
	Retries       uint16 `json:"retries" bson:"retries"`
	UpdatedAt     int32  `json:"updated_at" bson:"updated_at"`
	WriteTimout   uint32 `json:"write_timeout" bson:"write_timeout"`
}

type Latencies struct {
	Proxy   uint16 `json:"proxy" bson:"proxy"`
	Gateway uint16 `json:"kong" bson:"gateway"`
	Request uint16 `json:"request" bson:"request"`
}
