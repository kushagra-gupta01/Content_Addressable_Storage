package p2p

import(

)


//Peer is an interface that represents a remote node
type Peer interface{

}


//Transport is anything that handles communication
//between the nodes in the network. This can be of 
//the form(TCP, UDP, WebSockets,...) 
type Transport interface{
	ListenAndAccept() error
}

