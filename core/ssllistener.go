package core

import (
	"chimney3-go/privacy"
	"log"
	"net"
)

type MySSLListener interface {
	net.Listener
}

type SSLListenerImpl struct {
	RawListener   net.Listener
	Key           []byte
	II            privacy.EncryptThings
	ListenChannel chan MySSLSocket
}

func ListenSSL(host string, key []byte, i privacy.EncryptThings) (MySSLListener, error) {
	l, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}
	lss := &SSLListenerImpl{
		RawListener:   l,
		Key:           key,
		II:            i,
		ListenChannel: make(chan MySSLSocket),
	}

	go func() {
		for {
			conn, err := lss.RawListener.Accept()
			if err != nil {
				log.Println(" accept failed ", err)
				break
			}
			SetConnectTimeout(conn, 600)

			go func() {

				sock := NewMySSLSocket(conn, lss.II, lss.Key)
				err := sock.HandshakeServer()
				if err != nil {
					log.Println(" handshake failed ", err)
					sock.Close()
					return
				}
				log.Println(" handshake success ", conn.RemoteAddr().String())
				lss.ListenChannel <- sock
			}()
		}
	}()

	return lss, nil
}

func (l *SSLListenerImpl) Accept() (net.Conn, error) {

	conn, ok := <-l.ListenChannel
	if !ok {
		return nil, net.ErrClosed
	}
	log.Println(" accept success ", conn.RemoteAddr().String())
	return conn, nil
}

func (l *SSLListenerImpl) Close() error {
	close(l.ListenChannel)
	return l.RawListener.Close()
}

func (l *SSLListenerImpl) Addr() net.Addr {
	return l.RawListener.Addr()
}
