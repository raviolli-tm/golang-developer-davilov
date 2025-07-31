package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	tcpIn   *bufio.Reader
	tcpOut  *bufio.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{address: address, timeout: timeout, in: in, out: out}
}

func (s *TelnetClientImpl) Connect() error {
	if s.conn != nil {
		return nil
	}
	dial, err := net.DialTimeout("tcp", "localhost:4242", s.timeout)
	if err != nil {
		return err
	}
	s.conn = dial
	s.tcpIn = bufio.NewReader(s.conn)
	s.tcpOut = bufio.NewWriter(s.conn)

	return nil
}

func (s *TelnetClientImpl) Close() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}
	err = s.in.Close()
	if err != nil {
		return err
	}
	s.conn = nil
	s.tcpIn = nil
	s.tcpOut = nil
	return nil
}

func (s *TelnetClientImpl) Send() error {
	if s.conn == nil {
		return fmt.Errorf("no connection")
	}
	if s.in == nil {
		return fmt.Errorf("input stream is not set")
	}
	_, err := io.Copy(s.tcpOut, s.in)
	if err != nil {
		return err
	}
	return s.tcpOut.Flush()
}

func (s *TelnetClientImpl) Receive() error {
	if s.conn == nil {
		return fmt.Errorf("no connection")
	}
	if s.out == nil {
		return fmt.Errorf("output stream is not set")
	}
	_, err := io.Copy(s.out, s.tcpIn)
	return err
}
