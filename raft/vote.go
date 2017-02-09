package raft

//
// example RequestVote RPC arguments structure.
//
type RequestVoteArgs struct {
	// Your data here.
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

//
// example RequestVote RPC reply structure.
//
type RequestVoteReply struct {
	// Your data here.
	Term        int
	VoteGranted bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here.
	// do not need lock, send task to event loop across blockqueue Chan.
	ok := rf.deliver(&args, reply)
	if ok != nil {
		reply = nil
	}
}

func (rf *Raft) handleRequestVote(args *RequestVoteArgs) (RequestVoteReply, bool) {
	if args.Term < rf.CurrentTerm {
		//change the remote server which the vote come from to be a follower.
		return RequestVoteReply{rf.CurrentTerm, false}, false
	}

	if args.Term > rf.CurrentTerm {
		//change the local server to follower, reset CurrentTerm
		rf.updateCurrentTerm(args.Term, EmptyVote)
	} else if rf.VotedFor != EmptyVote && rf.VotedFor != args.CandidateId {
		// reject this remove server vote
		return RequestVoteReply{rf.CurrentTerm, false}, false
	}
	//log read/write lock
	rf.logmu.RLock()
	lastLogIndex := rf.LastIndex()
	lastLogTerm := rf.Log[lastLogIndex-rf.BaseIndex()].Term
	rf.logmu.RUnlock()
	if args.LastLogTerm > lastLogTerm || (args.LastLogTerm == lastLogTerm && args.LastLogIndex >= lastLogIndex) {
		//Vote
		rf.Vote(args.CandidateId)
		return RequestVoteReply{rf.CurrentTerm, true}, true
	} else {
		return RequestVoteReply{rf.CurrentTerm, false}, false
	}
	/*if lastLogIndex > args.LastLogIndex || lastLogTerm > args.LastLogTerm {
		// reject this remove server vote
		return RequestVoteReply{rf.CurrentTerm, false}, false
	} else {
		rf.VotedFor = args.CandidateId
		return RequestVoteReply{rf.CurrentTerm, true}, true
	}*/
}

func (rf *Raft) handleResponseVote(reply *RequestVoteReply) bool {

	if reply.Term > rf.CurrentTerm {
		rf.updateCurrentTerm(reply.Term, EmptyVote)
	}
	return reply.VoteGranted
}
