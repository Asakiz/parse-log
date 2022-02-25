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
	// FIX: verify the size of the type
	StartedAt uint32 `json:"started_at"`
}

type Request struct {
	Method string `json:"method" bson:"method"`
	URI    string `json:"uri" bson:"uri"`
	URL    string `json:"url" bson:"url"`
	// FIX: verify the size of the type
	Size        uint16        `json:"size" bson:"size"`
	Querystring interface{}   `json:"querystring" bson:"querystring"`
	Headers     RequestHeader `json:"headers" bson:"headers"`
}

type Response struct {
	// FIX: verify the size of the type
	Status uint16 `json:"status" bson:"status"`
	// FIX: verify the size of the type
	Size    uint16         `json:"size" bson:"size"`
	Headers ResponseHeader `json:"headers" bson:"headers"`
}

type AuthenticatedEntity struct {
	ConsumerID UUID `json:"consumer_id" bson:"consumer_id"`
}

type UUID struct {
	ID string `json:"uuid" bson:"uuid"`
}

type ID struct {
	ID string `json:"id" bson:"id"`
}

type Timestamp struct {
	// FIX: this type seems strange. ex of input: 1521555129
	CreatedAt int32 `json:"created_at" bson:"created_at"`
	// FIX: this type seems strange. ex of input: 1521555129
	UpdatedAt int32 `json:"updated_at" bson:"updated_at"`
}

type Route struct {
	Timestamp
	ID
	Hosts        string   `json:"hosts" bson:"hosts"`
	Methods      []string `json:"methods" bson:"methods"`
	Paths        []string `json:"paths" bson:"paths"`
	PreserveHost bool     `json:"preserve_host" bson:"preserve_host"`
	Protocols    []string `json:"protocols" bson:"protocols"`
	// FIX: verify the size of the type
	RegexPriority uint64       `json:"regex_priority" bson:"regex_priority"`
	Service       RouteService `json:"service" bson:"service"`
	StripPath     bool         `json:"strip_path" bson:"strip_path"`
}

type Service struct {
	Timestamp
	ID
	// FIX: verify the size of the type
	ConnectTimout uint32 `json:"connect_timeout" bson:"connect_timout"`
	Host          string `json:"host" bson:"host"`
	Name          string `json:"name" bson:"name"`
	Path          string `json:"path" bson:"path"`
	// FIX: verify the size of the type
	Port     uint16 `json:"port" bson:"port"`
	Protocol string `json:"protocol" bson:"protocol"`
	// FIX: the type could be not bigger enough
	ReadTimeout uint32 `json:"read_timeout" bson:"read_timeout"`
	// FIX: has a maximum number of tries?
	Retries uint16 `json:"retries" bson:"retries"`
	// FIX: the type could be not bigger enough
	WriteTimout uint32 `json:"write_timeout" bson:"write_timeout"`
}

type Latencies struct {
	// FIX: verify the size of the type
	Proxy uint16 `json:"proxy" bson:"proxy"`
	// FIX: verify the size of the type
	Gateway uint16 `json:"gateway" bson:"gateway"`
	// FIX: verify the size of the type
	Request uint16 `json:"request" bson:"request"`
}

type RequestHeader struct {
	Accept    string `json:"accept" bson:"accept"`
	Host      string `json:"host" bson:"host"`
	UserAgent string `json:"user-agent" bson:"user_agent"`
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

type RouteService struct {
	ID
}
