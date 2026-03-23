package raft

type State string

const (
	StateFollower  State = "follower"
	StateCandidate State = "candidate"
	StateLeader    State = "leader"
)

type Status struct {
	NodeID        string
	State         State
	ListenAddr    string
	TransportKind string
	PeerCount     int
	Started       bool
}
