package main

import "testing"
import "log"

var storage = NewStorage("/tmp")
var appid int64 = 0
var master *Master
func init() {
	master = NewMaster()
	master.Start()
}

func Test_Storage(t *testing.T) {
	im := &IMMessage{sender:1, receiver:2, content:"test"}
	msg := &Message{cmd:MSG_IM, body:im}
	msgid := storage.SaveMessage(msg)
	msg2 := storage.LoadMessage(msgid)
	if msg2 != nil {
		log.Println("msg2 cmd:", msg2.cmd)
	} else {
		log.Println("can't load msg:", msgid)
	}
}
 
func Test_Offline(t *testing.T) {
	im := &IMMessage{sender:1, receiver:2, content:"test"}
	msg := &Message{cmd:MSG_IM, body:im}
 
	for i := 0; i < 100; i++ {
		msgid := storage.SaveMessage(msg)
		storage.EnqueueOffline(msgid, appid, im.receiver)
		log.Println("enqueue msgid:", msgid)
	}
 
	log.Println("----------------------------")
	offs := storage.LoadOfflineMessage(appid, im.receiver)
	for _, emsg := range(offs) {
		log.Println("dequeue msgid:", emsg.msgid)
		storage.DequeueOffline(emsg.msgid, appid, im.receiver)
	}
}

func Test_Dequeue(t *testing.T) {
	im := &IMMessage{sender:1, receiver:2, content:"test"}
	msg := &Message{cmd:MSG_IM, body:im}


	msgid := storage.SaveMessage(msg)
	storage.EnqueueOffline(msgid, appid, im.receiver)
	log.Println("enqueue msgid:", msgid)
	storage.DequeueOffline(msgid, appid, im.receiver)
	storage.DequeueOffline(msgid, appid, im.receiver)
}

func Test_LoadLatest(t *testing.T) {
	im := &IMMessage{sender:1, receiver:2, content:"test"}
	msg := &Message{cmd:MSG_IM, body:im}

	msgid := storage.SaveMessage(msg)
	storage.EnqueueOffline(msgid, appid, im.receiver)

	im = &IMMessage{sender:1, receiver:2, content:"test2"}
	msg = &Message{cmd:MSG_IM, body:im}
	msgid = storage.SaveMessage(msg)
	storage.EnqueueOffline(msgid, appid, im.receiver)

	messages := storage.LoadLatestMessages(appid, im.receiver, 2)
	latest := messages[0]
	im2 := latest.msg.body.(*IMMessage)
	log.Println("sender:", im2.sender, " receiver:", im2.receiver, " content:", string(im2.content))


	latest = messages[1]
	im2 = latest.msg.body.(*IMMessage)
	log.Println("sender:", im2.sender, " receiver:", im2.receiver, " content:", string(im2.content))

}
