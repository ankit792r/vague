// Dispatch the incomming messages to handler
package backbone

import (
	"context"
	"fmt"
	"vagues/backbone/process"
	"vagues/backbone/session"
)

func (s *Server) dispatch(ctx context.Context, sess *session.Session, msg process.Message) {
	fmt.Println(msg.Params)
}
