package mqtt

import (
	"bytes"
	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
	"golang.org/x/net/context"
	"net"
)

type Config struct {
	Host string
	Port string
	User string
	Pass string
	CID  string
}

type MQTT struct {
	Config
	Messages chan Publish
	cl       *mqtt.ClientConn
}

type Publish struct {
	Topic string
	ID    uint16
	Msg   string
}

func New(cfg Config) *MQTT {
	return &MQTT{
		Config: cfg,
		Messages:   make(chan Publish),
	}
}

func (m *MQTT) Connect(ctx context.Context) error {
	conn, err := net.Dial("tcp", m.Host+":"+m.Port)

	if err != nil {
		panic(err)
	}

	m.cl = mqtt.NewClientConn(conn)
	m.cl.ClientId = m.CID

	err = m.cl.Connect(m.User, m.Pass)

	if err != nil {
		return err
	}

	go m.listen(ctx)

	return nil
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
Listen:
	for {
		select {
		case <-ctx.Done():
			break Listen
		case p := <-m.cl.Incoming:
			m.receive(p)
		}
	}
}

func (m *MQTT) receive(p *proto.Publish) {
	buf := new(bytes.Buffer)

	if err := p.Payload.WritePayload(buf); err != nil {
		panic(err)
	}

	pb := Publish{
		Topic: p.TopicName,
		ID:    p.MessageId,
		Msg:   buf.String(),
	}

	m.Messages <- pb
}
