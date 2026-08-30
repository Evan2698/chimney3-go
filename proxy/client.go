package proxy

import (
	"chimney3-go/core"
	"chimney3-go/privacy"
	"chimney3-go/utils"
	"net"
	"sync"
)

type proxyClient struct {
	Password     string
	I            privacy.EncryptThings
	LocalHost    string
	ProxyAddress string
	Exit         bool
	listener     net.Listener
	mu           sync.Mutex
}

// Client is the canonical proxy client name.
type Client = proxyClient

type ProxyClient interface {
	Serve() error
	Close() error
}

func (c *proxyClient) Serve() error {
	defer func() {
		utils.Recover("proxyClient.Serve")
	}()

	l, err := net.Listen("tcp", c.LocalHost)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.listener = l
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		if c.listener != nil {
			utils.CloseQuietly(c.listener)
			c.listener = nil
		}
		c.mu.Unlock()
	}()

	for {
		con, err := l.Accept()
		if err != nil {
			if c.Exit {
				return nil
			}
			return err
		}
		if c.Exit {
			utils.CloseQuietly(con)
			return nil
		}
		go c.serveOn(con)
	}
}

func (c *proxyClient) Close() error {
	c.mu.Lock()
	c.Exit = true
	listener := c.listener
	c.listener = nil
	c.mu.Unlock()
	utils.CloseAll(listener)
	return nil
}

func (c *proxyClient) serveOn(con net.Conn) {
	defer func() {
		utils.Recover("proxyClient.serveOn")
	}()

	defer utils.CloseQuietly(con)

	dst, err := c.connectToRemote()
	if err != nil {
		utils.LogError("proxyClient.serveOn handshake failed", err)
		return
	}
	defer utils.CloseQuietly(dst)

	waitForRelay := startBidirectionalRelay(con, dst)
	waitForRelay.Wait()
}

func (c *proxyClient) connectToRemote() (net.Conn, error) {
	key := privacy.MakeCompressKey(c.Password)
	dstIm, err := net.Dial("tcp", c.ProxyAddress)
	if err != nil {
		return nil, err
	}

	dst := core.NewMySSLSocket(dstIm, nil, key)
	if err := dst.HandshakeClient(); err != nil {
		utils.CloseQuietly(dst)
		return nil, err
	}
	if !dst.IsOk() {
		utils.CloseQuietly(dst)
		return nil, err
	}
	return dst, nil
}
