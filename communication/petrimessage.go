package communication

import (
  "bytes"
  "encoding/gob"
  "fmt"

  "github.com/FLAGlab/DCoPN/petrinet"
  "github.com/perlin-network/noise"
  "github.com/perlin-network/noise/log"
  "github.com/perlin-network/noise/payload"
  "github.com/pkg/errors"
)

// CommandType enums for PetriNodes communication
type CommandType string

const (
  // TransitionCommand to query transitions
  TransitionCommand  CommandType = "transitions"
  // MarksCommand to get number of marks on place
  MarksCommand        CommandType = "marks"
  // FireCommand to activate a fire event on the PetriNet
  FireCommand        CommandType = "fire"
  // PrintCommand to print the current state of the PetriNet
  PrintCommand       CommandType = "print"
  // VoteCommand to vote for a leader node
  VoteCommand        CommandType = "vote"
  // RequestVoteCommand to request votes for a leader node
  RequestVoteCommand CommandType = "requestvote"
  // AddToPlacesCommand to request mark addition (pos or neg) to a place
  AddToPlacesCommand CommandType = "addtoplaces"
  // RollBackTemporalPlacesCommand to request roll back of all temporal places
  RollBackTemporalPlacesCommand CommandType = "rollbacktemporal"
  // RollBackCommand to request roll back
  RollBackCommand CommandType = "rollback"
)

type petriMessage struct {
  Command CommandType
  Address string
  Term int
  FromType NodeType
  VoteGranted string
  Transitions []*petrinet.Transition
  RemoteTransitions map[int]*petrinet.RemoteTransition
  RemoteArcs []*petrinet.RemoteArc
  OpType petrinet.OperationType
  PetriContext string
  AskedPriority int
  imNew bool
  iveBeenFired bool
  SaveHistory bool
}

func (p petriMessage) String() string {
  return fmt.Sprintf("{Command: %v, Address: %v, Term: %v, FromType: %v, VoteGranted: %v, Transitions: %v, RemoteTransitions: %v, RemoteArcs: %v, OpType: %v, PetriContext: %v, AskedPriority: %v, imNew: %v}",
  p.Command, p.Address, p.Term, p.FromType, p.VoteGranted, p.Transitions, p.RemoteTransitions, p.RemoteArcs, p.OpType,
  p.PetriContext, p.AskedPriority, p.imNew)
}

func (petriMessage) Read(reader payload.Reader) (noise.Message, error) {
  byts, err := reader.ReadBytes()
  if err != nil {
    return nil, errors.Wrap(err, "failed to read msg")
  }
  var m petriMessage
  dec := gob.NewDecoder(bytes.NewReader(byts))
  if err := dec.Decode(&m); err != nil {
    return nil, errors.Wrap(err, "failed to decode msg")
  }

  return m, nil
}

func (m petriMessage) Write() []byte {
  var buf bytes.Buffer
  if err := gob.NewEncoder(&buf).Encode(m); err != nil {
    log.Info().Msgf("Got a fucking error: %v", err)
  }
  return payload.NewWriter(nil).WriteBytes(buf.Bytes()).Bytes()
}
