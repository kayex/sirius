package sirius

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kayex/sirius/mqtt"
	"github.com/kayex/sirius/slack"
)

type SyncAction string

const (
	NEW    SyncAction = "new"
	UPDATE            = "update"
	DELETE            = "delete"
)

type SyncMessage struct {
	Type SyncAction
	ID   slack.SecureID
}

type Sync interface {
	Sync(context.Context, *Service)
}

type SyncedService struct {
	service *Service
	sync    Sync
}

func (s *SyncedService) Start(ctx context.Context, u []User) {
	go s.service.Start(ctx, u)
	s.sync.Sync(ctx, s.service)
}

type MQTTSync struct {
	mqtt    *mqtt.MQTT
	topic   string
	rmt     *Remote
	service *Service
}

func (s *Service) WithSync(sync Sync) *SyncedService {
	return &SyncedService{
		service: s,
		sync:    sync,
	}
}

func NewMQTTSync(rmt *Remote, cfg mqtt.Config, topic string) *MQTTSync {
	return &MQTTSync{
		mqtt:  mqtt.New(cfg),
		topic: topic,
		rmt:   rmt,
	}
}

func (m *MQTTSync) Sync(ctx context.Context, s *Service) {
	m.service = s
	err := m.mqtt.Connect()

	if err != nil {
		panic(err)
	}

	m.mqtt.Subscribe(m.topic)
	m.start(ctx)
}

func (m *MQTTSync) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			m.mqtt.Disconnect()
			return
		case sm := <-m.mqtt.Messages:
			msg, err := parseSyncMessage(sm.Msg)

			if err != nil {
				continue
			}

			switch msg.Type {
			case UPDATE:
				m.service.DropUser(msg.ID)
				fallthrough
			case NEW:
				u, err := m.rmt.GetUser(msg.ID)

				if err != nil {
					break
				}

				m.service.AddUser(u)
			case DELETE:
				m.service.DropUser(msg.ID)
			}
		}

	}
}

func parseSyncMessage(msg string) (*SyncMessage, error) {
	split := strings.Split(msg, ":")

	if len(split) != 2 {
		return nil, errors.New(fmt.Sprintf("Invalid sync message %q", msg))
	}

	msgType := SyncAction(split[0])
	id := slack.SecureID{split[1]}

	switch msgType {
	case NEW, UPDATE, DELETE:
		return &SyncMessage{
			Type: msgType,
			ID:   id,
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown sync message type %q", msgType))
	}
}
