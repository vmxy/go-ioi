package main

import "github.com/vmxy/go-ioi/ioi"

type Session = ioi.Session

var Log = ioi.Log

func Listen(host string, port int, handle ioi.SessionHandle) {
	ioi.Listen(host, port, handle)
}

func Dail(host string, port int) (sess *ioi.Session, err error) {
	return ioi.Dail(host, port)
}
