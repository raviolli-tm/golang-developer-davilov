package main

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClientEdgeCases(t *testing.T) {
	t.Run("close without connect", func(t *testing.T) {
		client := NewTelnetClient("127.0.0.1:1", time.Second, io.NopCloser(&bytes.Buffer{}), io.Discard)
		require.NoError(t, client.Close())
		require.NoError(t, client.Close())
	})

	t.Run("send without connect", func(t *testing.T) {
		client := NewTelnetClient("127.0.0.1:1", time.Second, io.NopCloser(&bytes.Buffer{}), io.Discard)
		require.Error(t, client.Send())
	})

	t.Run("receive without connect", func(t *testing.T) {
		client := NewTelnetClient("127.0.0.1:1", time.Second, io.NopCloser(&bytes.Buffer{}), io.Discard)
		require.Error(t, client.Receive())
	})

	t.Run("connection error", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		address := listener.Addr().String()
		require.NoError(t, listener.Close())

		client := NewTelnetClient(address, 200*time.Millisecond, io.NopCloser(&bytes.Buffer{}), io.Discard)
		require.Error(t, client.Connect())
	})

	t.Run("receive eof from peer", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		defer listener.Close()

		client := NewTelnetClient(listener.Addr().String(), time.Second, io.NopCloser(&bytes.Buffer{}), io.Discard)
		require.NoError(t, client.Connect())
		defer client.Close()

		serverConn := make(chan net.Conn, 1)
		go func() {
			conn, acceptErr := listener.Accept()
			if acceptErr == nil {
				serverConn <- conn
			}
		}()

		conn := <-serverConn
		require.NoError(t, conn.Close())

		require.NoError(t, client.Receive())
	})
}
