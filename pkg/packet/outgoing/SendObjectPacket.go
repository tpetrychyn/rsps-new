package outgoing

import "github.com/tpetrychyn/rsps-comm-test/pkg/models"

type SendObjectPacket struct {
	Object *models.Actor
}

