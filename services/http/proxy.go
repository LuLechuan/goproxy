package http

type Proxy struct {
	ProxyName string
	IP        string
	IsDynamic bool
	User      string
	Pass      string
	Timestamp int
}

func (p *Proxy) SetIP(ip string) {
	p.IP = ip
}

func NewProxy(name string, ip string, isDynamic bool, user string, pass string, timestamp int) *Proxy {
	return &Proxy{
		ProxyName: name,
		IP:        ip,
		IsDynamic: isDynamic,
		User:      user,
		Pass:      pass,
		Timestamp: timestamp,
	}
}
