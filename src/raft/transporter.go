package raft

type Transporter interface {
	sendVote(from string, to string, request *VoteRequest) *VoteResponse
	sendAppend(from string, to string, request *AppendRequest) *AppendResponse
}

// the vote channel and append channel of the server
type Channel struct {
	voteChan   chan *Message
	appendChan chan *Message
}

func newChannel(voteChan chan *Message, appendChan chan *Message) *Channel {
	return &Channel{voteChan, appendChan}
}

type MemTransporter struct {
	channels map[string]*Channel
}

func newMemTransporter(channels map[string]*Channel) Transporter {
	return &MemTransporter{channels}
}

func (trans *MemTransporter) sendVote(from string, to string, request *VoteRequest) *VoteResponse {
	responseChan := make(chan interface{}, 1)
	message := &Message{from, to, request, responseChan}
	trans.channels[to].voteChan <- message
	resp := <-responseChan
	return resp.(*VoteResponse)

}

func (trans *MemTransporter) sendAppend(from string, to string, request *AppendRequest) *AppendResponse {
	responseChan := make(chan interface{}, 1)
	message := &Message{from, to, request, responseChan}
	trans.channels[to].appendChan <- message
	resp := <-responseChan
	return resp.(*AppendResponse)
}
