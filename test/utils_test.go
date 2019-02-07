package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sumwhere/chat"
	"testing"
	"time"

	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	cass "github.com/gocql/gocql"
	"google.golang.org/api/option"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("토픽: %s\n", msg.Topic())
	fmt.Printf("메시지: %s\n", msg.Payload())
}

func TestMqtt(t *testing.T) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.1.11:1883")
	opts.SetClientID("go-simple")
	opts.SetDefaultPublishHandler(f)
	opts.Username = "qkrqjadn"
	opts.Password = "1q2w3e4r"

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	token := c.Connect()
	token.Wait()

	assert.Nil(t, token.Error(), token.Error())

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe("chat/*/*", 2, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 5; i++ {
		t.Log("chat/1/1", i)
		text := fmt.Sprintln("유저들", i)
		token := c.Publish("chat/1/2", 2, false, text)
		token.Wait()
	}

	for i := 0; i < 5; i++ {
		t.Log("chat/2/1", i)
		text := fmt.Sprintln("유저들", i)
		token := c.Publish("chat/2/2", 2, false, text)
		token.Wait()
	}

	for i := 0; i < 5; i++ {
		t.Log("chat/1", i)
		text := fmt.Sprintln("채팅방들", i)
		token := c.Publish("chat/1", 2, false, text)
		token.Wait()
	}

	for i := 0; i < 5; i++ {
		t.Log("chat/2", i)
		text := fmt.Sprintln("채팅방들", i)
		token := c.Publish("chat/2", 2, false, text)
		token.Wait()
	}

	time.Sleep(3 * time.Second)

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("chat/*"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}

func TestCassandra(t *testing.T) {
	cluster := cass.NewCluster("210.100.238.118")
	cluster.Port = 59042
	cluster.Keyspace = "zip"
	cluster.Consistency = cass.Quorum
	session, err := cluster.CreateSession()

	defer session.Close()

	err = session.Query(`CREATE TABLE IF NOT EXISTS room_1 (id timeuuid, dataType int, data blob , PRIMARY KEY(id))`).Exec()
	assert.Nil(t, err, err.Error())

	for i := 0; i < 10; i++ {
		err = session.Query(`INSERT INTO room_1 (id, data, dataType) VALUES (?, ?, ?)`, cass.TimeUUID(), "hello world", i).Exec()
		assert.Nil(t, err, err.Error())
	}

	var id cass.UUID
	var data []byte
	var dataType int

	sampleUUID := "7d41ebf3-b9bc-11e8-9adf-b8e8563ccf5c"

	iter := session.Query(`SELECT id, data, dataType FROM room_1 WHERE id > ?  ALLOW FILTERING`, sampleUUID).Iter()

	for iter.Scan(&id, &data, &dataType) {
		t.Log(id.Time(), data, dataType)
	}

	err = iter.Close()
	assert.Nil(t, err, err.Error())
}

func TestCloudMessage(t *testing.T) {
	opt := option.WithCredentialsFile("./sumwhere-8f900-firebase-adminsdk-zhjsl-f9265d989f.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	assert.Nil(t, err, err.Error())

	client, err := app.Messaging(context.Background())
	assert.Nil(t, err, err.Error())

	registrationToken := "fVuRI0OAQfo:APA91bETAj86O4KsvOX3m0R9ULOfAn3R5ord9tGFSL5Hhihgsa1NgM_gZ93W-4nBafsvpRkjgsY6o3d9kgoPxtrIT5OTn4O6Nt1Kr53vP4RNT_8F3_P14xoLnvxcqjVRPJ4bmUKbLnnX"
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}
	for i := 0; i < 10; i++ {
		response, err := client.Send(context.Background(), message)
		assert.Nil(t, err, err.Error())
		t.Log("success", response)
	}
}

func Test_Mymqtt(t *testing.T) {

	c := make(chan int64)
	mqtt := chat.Mqtt{}.New("ems01", "1883")
	err := mqtt.NewSubscribe("chat/#", func(client MQTT.Client, message MQTT.Message) {
		fmt.Println(message.Topic())
	})

	assert.Nil(t, err, err.Error())

	//err = mqtt.NewSubscribe("chat/2/#", func(client MQTT.Client, message MQTT.Message) {
	//	fmt.Println(message.Topic())
	//})
	//Ok(t, err)

	for i := 1; i < 10; i++ {
		err = mqtt.NewPublish("chat/1", "1")
		assert.Nil(t, err, err.Error())
	}

	for i := 1; i < 10; i++ {
		err = mqtt.NewPublish("go-mqtt/2", "2/1")
		assert.Nil(t, err, err.Error())
	}

	for i := 1; i < 10; i++ {
		err = mqtt.NewPublish("chat/3/3", "3/3")
		assert.Nil(t, err, err.Error())
	}

	for i := 1; i < 10; i++ {
		err = mqtt.NewPublish("hi/hello", "3/3")
		assert.Nil(t, err, err.Error())
	}

	defer func() {
		mqtt.Close()
		mqtt.UnSubscribe("*")
	}()
	time.Sleep(3 * time.Second)
	<-c
}
