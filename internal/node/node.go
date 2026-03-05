// Package node
// * listen server for messeges
package node

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/smirnovalles/yt/internal/protocol"
)

type Node struct {
	id         string
	port       int
	address    string
	peerStatus PeerStatus
	nodeStatus NodeStatus
	listener   net.Listener

	connections map[string]net.Conn
	connMutex   sync.RWMutex

	stop chan struct{}
}

type PeerStatus int

const (
	PeerStatusDisconnected PeerStatus = iota
	PeerStatusConnecting
	PeerStatusConnected
)

type NodeStatus int

const (
	NodeStatusStop  NodeStatus = iota //0 stop
	NodeStatusStart                   //1 start
)

func New(id string, ip string, port int) *Node {
	n := Node{
		id:          id,
		address:     ip,
		port:        port,
		nodeStatus:  NodeStatusStop,
		peerStatus:  PeerStatusDisconnected,
		connections: make(map[string]net.Conn),
	}

	return &n
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %s %s:%d", n.id, n.address, n.port)
}

func (n *Node) Start() error {

	addr := fmt.Sprintf("%s:%d", n.address, n.port)
	listener, err := net.Listen("tcp", addr)
	//открыт закрываем в Stop

	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	n.listener = listener
	n.nodeStatus = NodeStatusStart
	n.stop = make(chan struct{})

	go n.acceptLoop()

	return nil
}

func (n *Node) Stop() (err error) {

	if n.nodeStatus == NodeStatusStop {
		return nil
	}

	close(n.stop)

	//delete conn
	for address := range n.connections {
		n.deleteConn(address)
	}

	if n.listener != nil {
		err = n.listener.Close()
		// log err
	}

	n.nodeStatus = NodeStatusStop
	n.peerStatus = PeerStatusDisconnected

	return nil
}

func (n *Node) acceptLoop() {

	fmt.Println("Accept loop starting ... ")

	for {
		select {
		case <-n.stop:
			fmt.Println("Stopping accept loop")
			return
		default:
		}

		conn, err := n.listener.Accept()
		if err != nil {
			select {
			case <-n.stop:
				return
			default:
				fmt.Printf("Accept error: %v\n", err)
			}
			continue
		}

		go n.handleClient(conn, false)
	}
}

func (n *Node) saveConn(addr string, conn net.Conn) {

	n.connMutex.Lock()
	n.connections[addr] = conn
	n.connMutex.Unlock()

}

func (n *Node) loadConn(addr string) net.Conn {
	n.connMutex.RLock()
	conn, ok := n.connections[addr]
	n.connMutex.RUnlock()
	if ok {
		return conn
	}
	return nil
}

func (n *Node) deleteConn(addr string) {

	n.connMutex.Lock()
	conn, ok := n.connections[addr]
	delete(n.connections, addr)
	n.connMutex.Unlock()
	if conn != nil && ok {
		conn.Close()
	}

}

func (n *Node) handleClient(conn net.Conn, isOutgoing bool) {

	addr := conn.RemoteAddr().String()

	if isOutgoing {
		fmt.Printf("Started handler for outgoing connection to %s\n",
			addr)
	} else {

		n.saveConn(addr, conn)

		fmt.Printf("New incoming connection from %s\n",
			addr)
	}

	defer n.deleteConn(addr)

	reader := bufio.NewReader(conn)

	for {

		//read header
		header := make([]byte, protocol.MessageHeaderLen)

		_, err := io.ReadFull(reader, header)

		if err != nil {
			break
		}

		msgType, msgLen, err := protocol.GetHeader(header)

		if err != nil {
			break
		}

		payload := make([]byte, msgLen)

		_, err = io.ReadFull(reader, payload)
		if err != nil {
			break
		}

		m, err := protocol.Decode(msgType, payload)

		switch m.Type() {
		case protocol.TypeHandShake:
			n.incomingHandShake(m)
		case protocol.TypeTextMessage:
			n.incomingText(m)
		default:
			//TODO
		}

		// TODO: отправка ответа
		// conn.Write([]byte("OK"))
	}
}

func (n *Node) incomingHandShake(m protocol.Message) {
	fmt.Printf("Get handshake:%s", m.(*protocol.HandShakeMessage).GetID())

}

func (n *Node) incomingText(m protocol.Message) {
	fmt.Printf("Get text:%s", m.(*protocol.TextMessage).GetText())
}

func (n *Node) Send(peerAddr string, data []byte) error {

	conn, err := n.Connect(peerAddr)

	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err

}

func (n *Node) Connect(addr string) (conn net.Conn, err error) {

	conn = n.loadConn(addr)

	if conn != nil {
		return conn, nil
	}

	conn, err = net.Dial("tcp", addr)

	if err != nil {
		return nil, fmt.Errorf("connect failed: %w", err)
	}

	n.saveConn(addr, conn)

	fmt.Printf("Node %s connected to %s\n", n.id, addr)
	go n.handleClient(conn, true)

	return conn, nil
}
