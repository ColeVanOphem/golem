package p2p

// HandshakeFunc is ...?
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
