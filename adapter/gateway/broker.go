package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Go-routine-4995/ossfrontend/domain"
	"github.com/nats-io/nats.go"
	"time"
)

const (
	messageGet = iota
	messageGetPaged
	messageCreate
	messageDelete

	subjectRouter  = "ns.oss.router"
	subjectAccount = "ns.oos.account"
	timeout        = 5000    // ms
	maxMessageSize = 8000000 // bytes
)

type message struct {
	Mtype int    `json:"mtype"`
	Data  []byte `json:"data"`
}

type Broker struct {
	urlBroker string
	con       *nats.Conn
}

func NewBroker(u string) *Broker {
	c, err := connect(u)

	if err != nil {
		fmt.Println("Broker connection error: ", err)
	}
	return &Broker{
		urlBroker: u,
		con:       c,
	}
}

func connect(u string) (*nats.Conn, error) {
	nc, err := nats.Connect(u)

	return nc, err
}

func (b *Broker) CreateRoutersRequest(r []domain.Router, tenant string) ([]byte, error) {
	var (
		d   []byte
		err error
		msg *nats.Msg
		m   message
	)
	d, err = json.Marshal(r)
	if err != nil {
		return nil, err
	}

	if len(d) > maxMessageSize {
		return nil, errors.New("message size too large")
	}

	m.Mtype = messageCreate
	m.Data = d
	d, err = json.Marshal(m)
	if err != nil {
		return nil, err
	}

	msg, err = b.con.Request(subjectRouter, d, time.Duration(time.Millisecond*timeout))
	if err != nil {
		return nil, err
	}
	fmt.Println(string(msg.Data))

	return msg.Data, nil
}

func (b *Broker) DeleteRoutersRequest(r []domain.Router, tenant string) error {
	var (
		d   []byte
		err error
		msg *nats.Msg
		m   message
	)
	d, err = json.Marshal(r)
	if err != nil {
		return err
	}
	if len(d) > maxMessageSize {
		return errors.New("message size too large")
	}

	m.Mtype = messageDelete
	m.Data = d
	d, err = json.Marshal(m)
	if err != nil {
		return err
	}

	msg, err = b.con.Request(subjectRouter, d, time.Duration(time.Millisecond*timeout))
	if err != nil {
		return err
	}

	fmt.Println(string(msg.Data))

	return nil
}

func (b *Broker) GetRoutersPage(paginationByte []byte, tenant string) ([]domain.Router, error) {
	var (
		err     error
		routers []domain.Router
		msg     *nats.Msg
		m       message
		d       []byte
	)

	m.Mtype = messageGetPaged
	m.Data = paginationByte
	d, err = json.Marshal(m)
	if err != nil {
		return nil, err
	}

	msg, err = b.con.Request(subjectRouter, d, time.Duration(time.Millisecond*timeout))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(msg.Data, &routers)
	if err != nil {
		return nil, err
	}

	return routers, nil
}

func (b *Broker) GetRouters(r domain.Router, tenant string) (domain.Router, error) {
	var (
		err error
		ret domain.Router
		msg *nats.Msg
		d   []byte
		m   message
	)

	d, err = json.Marshal(r)
	if err != nil {
		return ret, err
	}
	if len(d) > maxMessageSize {
		return ret, errors.New("message size too large")
	}

	m.Mtype = messageGet
	m.Data = d
	d, err = json.Marshal(m)
	if err != nil {
		return ret, err
	}

	msg, err = b.con.Request(subjectRouter, d, time.Duration(time.Millisecond*timeout))
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(msg.Data, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
