package mqtt

import (
	"bytes"
	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
	"golang.org/x/net/context"
	"net"
)

type MQTTConfig struct {
	Host string
	Port string
	User string
	Pass string
	CID  string
}

type MQTT struct {
	MQTTConfig
	Messages chan Publish
	cl       *mqtt.ClientConn
}

type Publish struct {
	Topic string
	ID    uint16
	Msg   string
}

func New(cfg MQTTConfig) *MQTT {
	return &MQTT{
		MQTTConfig: cfg,
		Messages:   make(chan Publish),
	}
}

func (m *MQTT) Connect(ctx context.Context) {
	conn, err := net.Dial("tcp", m.Host+":"+m.Port)

	if err != nil {
		panic(err)
	}

	m.cl = mqtt.NewClientConn(conn)
	m.cl.ClientId = m.CID

	m.cl.Connect(m.User, m.Pass)

	go m.listen(ctx)
}

func (m *MQTT) Subscribe(topic string) {

	sub := []proto.TopicQos{
		{
			Topic: topic,
			Qos:   proto.QosAtLeastOnce,
		},
	}

	m.cl.Subscribe(sub)
}

func (m *MQTT) listen(ctx context.Context) {
	defer m.cl.Disconnect()
	defer close(m.Messages)

Listen:
	for {
		select {
		case p := <-m.cl.Incoming:
			m.receive(p)
		case <-ctx.Done():
			break Listen
		}
	}
}

func (m *MQTT) receive(p *proto.Publish) {
	buf := new(bytes.Buffer)

	if err := p.Payload.WritePayload(buf); err != nil {
		panic(err)
	}

	s := buf.String()

	pb := Publish{
		Topic: p.TopicName,
		ID:    p.MessageId,
		Msg:   s,
	}

	m.Messages <- pb
}
