package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttCli mqtt.Client
var deviceName string
var deviceID int = 0

func main() {
	mqttCli = creatMQTT("1124")
	fmt.Printf("--- 创建新的mqtt连接\n")

	//随机设备名
	rand.Seed(time.Now().Unix())
	deviceName = fmt.Sprintf("Schneider PLC%v", rand.Intn(10000))

	pub("register", deviceName)
	//time.Sleep(time.Second)

	sub(deviceName+"_id", messagePubHandler)
	sub(deviceName+"_inquiry", messagePubHandler)
	sub(deviceName+"_broadcast", messagePubHandler)
	sub(deviceName+"_delay", messagePubHandler)
	sub(deviceName+"_result", messagePubHandler)


	for i := 0; i <= 20; i++ {
		fmt.Printf("--- %v\n", deviceID)
		time.Sleep(time.Second)
	}

}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//var aa string = "11223"

	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	recive := fmt.Sprintf("%s", msg.Payload())
	if recive == "0" {
		fmt.Printf("--- %v\n", "设备名重复")
	}
	reciveTopic := fmt.Sprintf("%s", msg.Topic())
	if reciveTopic == deviceName+"_id" {
		deviceID, _ = strconv.Atoi(recive)
		fmt.Printf("--- %v\n", deviceID)
	} else if reciveTopic == deviceName+"_inquiry" {
		time.Sleep(time.Second)
		deviceIDStr := fmt.Sprintf("%v", deviceID)
		pub(deviceName+"_inquiry_send", deviceIDStr)
	} else if reciveTopic == deviceName+"_broadcast" {
		time.Sleep(time.Second)
		boardIDStr := fmt.Sprintf("%s", msg.Payload())
		pub(deviceName+"_broadcast_send", boardIDStr)
	} else if reciveTopic == deviceName+"_delay" {
		pub(deviceName+"_delay_send", recive)
	} else if reciveTopic == deviceName+"_result" {
		resultStr := fmt.Sprintf("%s", msg.Payload())
		if resultStr == "SUCCESS" {
			fmt.Printf("+----------------+\n")
			fmt.Printf("|测试结果：SUCESS|\n")
			fmt.Printf("+----------------+\n")
		} else {
			fmt.Printf("+----------------+\n")
			fmt.Printf("|测试结果：FAIL  |\n")
			fmt.Printf("+----------------+\n")
		}
	} else if reciveTopic == deviceName+"_alarmt" {
                pub(deviceName+"_alarm", "alarm")
	}

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func creatMQTT(connectName string) mqtt.Client {
	var broker = "0.0.0.0"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(connectName)
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//fmt.Printf("--------14122-0  %+v\n", client)

	//sub(client)
	//publish(client)

	//client.Disconnect(250)

	return client
}

//Pub 发送
func pub(topic string, message string) {

	// fmt.Printf("--------14122-1  %+v\n", client)

	//num := 10
	//fmt.Println("--------14122  *********************************************************************************************************")
	//for i := 0; i < num; i++ {
	//fmt.Println("--------14122  ", i)
	//text := fmt.Sprintf("Message %d", i)
	token := mqttCli.Publish(topic, 0, false, message)
	token.Wait()
	fmt.Printf("【发送mqtt消息】topic：%v  massage：%v\n", topic, message)
	//time.Sleep(time.Second)
	//}
}

//Sub 订阅
func sub(topic string, rtu func(client mqtt.Client, msg mqtt.Message)) {

	fmt.Printf("【订阅topic】：  %+v\n", topic)

	//topic := "11223"
	token := mqttCli.Subscribe(topic, 0x00, messagePubHandler)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
