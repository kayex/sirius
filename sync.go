package sirius

import (
	"github.com/kayex/sirius/mqtt"
	"golang.org/x/net/context"
	"strings"
)

type SyncAction string

const (
	NEW    SyncAction = "new"
	UPDATE            = "update"
	DELETE            = "delete"
)

type SyncMessage struct {
	Type  SyncAction
	Token string
}

type Sync interface {
	Sync(s *Service)
}

type SyncedService struct {
	service *Service
	sync    Sync
}

func (s *SyncedService) Start(ctx context.Context, u []User) {
	go s.service.Start(ctx, u)
	s.sync.Sync(s.service)
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

func (m *MQTTSync) Sync(s *Service) {
	m.service = s
	m.mqtt.Connect(context.TODO())
	m.mqtt.Subscribe(m.topic)

	m.start()
}

func (m *MQTTSync) start() {
	for msg := range m.mqtt.Messages {
		msg, ok := parseSyncMessage(msg.Msg)

		if !ok {
			continue
		}

		switch msg.Type {
		case UPDATE:
			m.service.DropUserWithToken(msg.Token)
			fallthrough
		case NEW:
			u, err := m.rmt.GetUser(msg.Token)

			if err != nil {
				break
			}

			m.service.AddUser(u)
		case DELETE:
			m.service.DropUserWithToken(msg.Token)
		}
	}
}

func parseSyncMessage(msg string) (*SyncMessage, bool) {
	split := strings.Split(msg, ":")

	if len(split) != 2 {
		return nil, false
	}

	msgType := SyncAction(split[0])
	token := split[1]

	switch msgType {
	case NEW:
		fallthrough
	case UPDATE:
		fallthrough
	case DELETE:
		return &SyncMessage{
			Type:  msgType,
			Token: token,
		}, true
	default:
		return nil, false
	}
}