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
	ControlAddr   string
	RaftAddr      string
	TransportKind string
	LeaderID      string
	PeerCount     int
	Started       bool
}
