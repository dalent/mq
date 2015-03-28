package rabbit

import (
	"errors"
	"github.com/dalent/mq"
	"github.com/streadway/amqp"
)

//我暂时只想封装下rabbit的接口
type rabbitAdapter struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue
}

var rabbit rabbitAdapter

func init() {
	mq.Register(mq.RABBIT, NewRabbitAdapter)
}

func NewRabbitAdapter() mq.Provider {
	return &rabbitAdapter{}
}
func (p *rabbitAdapter) SetUp(url, queue string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	p.conn = conn

	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}

	p.ch = ch

	q, err := ch.QueueDeclare(queue, //name
		true,  //durable
		false, //deleted when unused
		false, //exclusive
		false, //no-wait
		nil)   //argument
	if err != nil {
		return err
	}

	p.queue = &q
	return nil
}

func (p *rabbitAdapter) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
	p.conn = nil
	p.ch = nil
	p.queue = nil
}
func (p *rabbitAdapter) Check() error {
	if p.conn == nil || p.ch == nil || p.queue == nil {
		return errors.New("rabbit provider not set up")
	}

	return nil
}

func (p *rabbitAdapter) NewWork(exchange string, bytes []byte) error {
	err := p.Check()
	if err != nil {
		return err
	}

	return p.ch.Publish(
		exchange,     //exchange
		p.queue.Name, // routing key
		false,        //mandatory
		false,        //immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/json",
			Body:         bytes,
		})
}

func (p *rabbitAdapter) Worker(qos int, action mq.Action) error {
	err := p.Check()
	if err != nil {
		return err
	}

	err = p.ch.Qos(qos, //prefetch size
		0, //prefetch size
		false,
	) //global

	if err != nil {
		return err
	}

	msgs, err := p.ch.Consume(
		p.queue.Name, //queue
		"",           //consume
		true,         //auto-ack
		false,        //exclusive
		false,        //no-local
		false,        //no wait
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		action.Do(msg.Body)
	}

	return nil
}

func (p *rabbitAdapter) Publish(exchange string, bytes []byte) error {
	err := p.Check()
	if err != nil {
		return err
	}

	return p.ch.Publish(
		exchange,     //exchange
		p.queue.Name, // routing key
		false,        //mandatory
		false,        //immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        bytes,
		})
}
func (p *rabbitAdapter) Consume(action mq.Action) error {
	err := p.Check()
	if err != nil {
		return err
	}

	msgs, err := p.ch.Consume(
		p.queue.Name, //queue
		"",           //consume
		true,         //auto-ack
		false,        //exclusive
		false,        //no-local
		false,        //no wait
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		action.Do(msg.Body)
	}

	return nil
}
