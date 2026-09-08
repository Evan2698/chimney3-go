package core

import (
	"chimney3-go/privacy"
	"log"
	"net"
	"sync"
)

type MySSLListener interface {
	net.Listener
}

type SSLListenerImpl struct {
	RawListener   net.Listener
	Key           []byte
	II            privacy.EncryptThings
	ListenChannel chan MySSLSocket
	closeOnce     sync.Once
	done          chan struct{}
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
		done:          make(chan struct{}),
	}

	go func() {
		for {
			conn, err := lss.RawListener.Accept()
			if err != nil {
				log.Println(" accept failed ", err)
				break
			}
			if err := SetConnectTimeout(conn, 600); err != nil {
				log.Println(" set handshake timeout failed ", err)
				_ = conn.Close()
				continue
			}

			go func() {

				sock := NewMySSLSocket(conn, lss.II, lss.Key)
				err := sock.HandshakeServer()
				if err != nil {
					log.Println(" handshake failed ", err)
					sock.Close()
					return
				}
				if err := SetConnectTimeout(conn, 0); err != nil {
					log.Println(" clear handshake timeout failed ", err)
					sock.Close()
					return
				}
				log.Println(" handshake success ", conn.RemoteAddr().String())
				select {
				case lss.ListenChannel <- sock:
				case <-lss.done:
					_ = sock.Close()
				}
			}()
		}
	}()

	return lss, nil
}

func (l *SSLListenerImpl) Accept() (net.Conn, error) {
	if l == nil || l.ListenChannel == nil {
		return nil, net.ErrClosed
	}
	if l.done == nil {
		conn := <-l.ListenChannel
		return conn, nil
	}
	select {
	case conn := <-l.ListenChannel:
		return conn, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

func (l *SSLListenerImpl) Close() error {
	if l == nil {
		return nil
	}
	l.closeOnce.Do(func() {
		if l.done != nil {
			close(l.done)
		}
	})
	if l.RawListener != nil {
		return l.RawListener.Close()
	}
	return nil
}

func (l *SSLListenerImpl) Addr() net.Addr {
	return l.RawListener.Addr()
}
