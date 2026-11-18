// graft:synth:begin {"record":"prov-2026-69661708","version":1,"gen":"sha256:9c11891d45f78ed26ba4bf0b43a67e85a5a15a8cdbe44525c7f2b061129a4fa7","tier":"low","body":"f0b1022316fe640b85aeac588d236d2a59407752fc5ce36fb5a11006ecd42a3a"}
package main

import "net"

func (s *Server) runningServeBody(conn net.Conn) {
	panic("graft: not converged")
}

func (c *Client) runningSendBody(line string) {
	panic("graft: not converged")
}

func (s *Server) runningReplyBody(c *Client, code string, text string) {
	panic("graft: not converged")
}

func (c *Client) runningPrefixBody() string {
	panic("graft: not converged")
}

func (s *Server) runningDisconnectBody(c *Client, reason string) {
	panic("graft: not converged")
}

// graft:synth:end {"record":"prov-2026-69661708"}
