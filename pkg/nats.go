package pkg

import (
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"time"
)

func InitNats(inProcess bool, enableLogging bool) (*nats.Conn, nats.JetStreamContext, *server.Server, error) {

	opts := &server.Options{
		ServerName:      "embedded_server",
		DontListen:      inProcess,
		JetStream:       true,
		JetStreamDomain: "embedded",
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, nil, nil, err
	}

	if enableLogging {
		ns.ConfigureLogger()
	}
	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, nil, nil, err
	}

	clientOpts := []nats.Option{}
	if inProcess {
		clientOpts = append(clientOpts, nats.InProcessServer(ns))
	}

	nc, err := nats.Connect(nats.DefaultURL, clientOpts...)
	if err != nil {
		return nil, nil, nil, err
	}

	streamConfig := nats.StreamConfig{
		Name:        "messages",
		Description: "",
		Subjects:    []string{"messages"},
		Retention:   nats.LimitsPolicy,
		MaxAge:      24 * time.Hour,
		Replicas:    1,
		NoAck:       false,
	}
	js, err := nc.JetStream()
	if err != nil {
		return nil, nil, nil, err
	}

	_, err = js.AddStream(&streamConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	return nc, js, ns, err
}
