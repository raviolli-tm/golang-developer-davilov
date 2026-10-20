package main

import (
	"bufio"
	"errors"
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

	conn, err := net.DialTimeout("tcp", s.address, s.timeout)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", s.address, err)
	}

	s.conn = conn
	s.tcpIn = bufio.NewReader(conn)
	s.tcpOut = bufio.NewWriter(conn)

	return nil
}

func (s *TelnetClientImpl) Close() error {
	var result error

	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			result = err
		}
		s.conn = nil
	}

	if s.in != nil {
		if err := s.in.Close(); err != nil && result == nil {
			result = err
		}
	}

	s.tcpIn = nil
	s.tcpOut = nil

	return result
}

func (s *TelnetClientImpl) Send() error {
	if s.conn == nil || s.tcpOut == nil {
		return errors.New("no connection")
	}
	if s.in == nil {
		return errors.New("input stream is not set")
	}

	if _, err := io.Copy(s.tcpOut, s.in); err != nil {
		return fmt.Errorf("send data: %w", err)
	}

	if err := s.tcpOut.Flush(); err != nil {
		return fmt.Errorf("flush data: %w", err)
	}

	return nil
}

func (s *TelnetClientImpl) Receive() error {
	if s.conn == nil || s.tcpIn == nil {
		return errors.New("no connection")
	}
	if s.out == nil {
		return errors.New("output stream is not set")
	}

	if _, err := io.Copy(s.out, s.tcpIn); err != nil {
		return fmt.Errorf("receive data: %w", err)
	}

	return nil
}
