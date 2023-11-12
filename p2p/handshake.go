package p2p

//HandshakeFunc....
type HandshakeFunc func (a any) error

func NOPHandshakeFunc(any) error {return nil}