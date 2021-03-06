package forwarder

import (
	"net"
	"time"

	l "github.com/Sirupsen/logrus"
	"github.com/mikioh/tcp"
)

func init() {
	RegisterForwarder("TCP", newTCP)
}

// TCP forwarder
type TCP struct {
	BaseForwarder
	conn   *net.TCPConn
	server string
	port   string
}

// newTCP returns a new TCP forwarder
func newTCP(
	initialBufferSize int,
	log *l.Entry) Forwarder {

	t := new(TCP)
	t.name = "TCP"

	t.maxBufferSize = initialBufferSize
	t.log = log
	return t
}

// Configure the TCP forwarder
func (t *TCP) Configure(configMap map[string]interface{}) {
	if server, exists := configMap["server"]; exists {
		t.server = server.(string)
	} else {
		t.log.Error("There was no server specified, there won't be any emissions")
	}

	if port, exists := configMap["port"]; exists {
		t.port = port.(string)
	} else {
		t.log.Error("There was no port specified , there won't be any emissions")
	}
	t.configureCommonParams(configMap)
}

// Run runs the forwarder main loop
func (t *TCP) Run() {
	var c *tcp.Conn
	var conn net.Conn
	var err error
	conn, err = net.Dial("tcp", t.server+":"+t.port)
	if err != nil {
		t.log.Error("Could not connect to remote TCP host")
		return
	}

	c, err = tcp.NewConn(conn)
	if err != nil {
		t.log.Error("Could not cretae TCPConn")
		conn.Close()
		return
	}

	c.Cork()
	t.conn = &c.TCPConn
	t.conn.SetKeepAlive(true)
	t.conn.SetKeepAlivePeriod(time.Duration(t.KeepAliveInterval()) * time.Second)

	t.run(t.emitMsg)
}

func (t *TCP) emitMsg(m []byte) bool {

	_, err := t.conn.Write(m)
	if err != nil {
		t.log.Error("Failed to send message to TCP endpoint")
		return false
	}

	return true
}
