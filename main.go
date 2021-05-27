package main

import (
	"encoding/json"
	. "github.com/lore-lml/go_bioenpro4to_channel_manager"
)

type Message struct {
	ActorId    string `json:"actor_id"`
	Msg        string `json:"msg"`
	Visibility string `json:"visibility"`
}

func newMessage(actorId string, public bool) Message {
	visibility := "PUBLIC"
	if !public {
		visibility = "PRIVATE"
	}
	return Message{
		ActorId:    actorId,
		Msg:        "This is a message",
		Visibility: visibility,
	}
}

func testCreateChannelTree(statePsw string, mainnet bool, keyNonce *KeyNonce) ChannelInfo {
	root := NewRootChannel(mainnet)
	defer root.Drop()

	info := root.Open(statePsw) //Apre ed inizializza i channel dei primi due layer dell' albero, essendo fissi

	//Prova di creazione dei layer 3 e 4 per lo specifico actor di una determinata categoria in una certa data
	root.GetCreateDailyChannelManager(Trucks, "XASD", statePsw, 25, 5, 2021)
	root.GetCreateDailyChannelManager(Trucks, "XASD", statePsw, 26, 5, 2021)
	root.GetCreateDailyChannelManager(Trucks, "XASD2", statePsw, 25, 5, 2021)
	root.GetCreateDailyChannelManager(Scales, "SCALE1", statePsw, 28, 5, 2021)

	//Crea un daily channel e si prova a inviare un messaggio al suo  interno
	dailyCh := root.GetCreateDailyChannelManager(Trucks, "XASD", statePsw, 25, 5, 2021)
	defer dailyCh.Drop()
	//Creazione del messaggio Serializzando la struttura Message in un json byte array
	public, _ := json.Marshal(newMessage("XASD", true))
	private, _ := json.Marshal(newMessage("XASD", false))
	packet := NewRawPacket(public, private)
	defer packet.Drop()
	dailyCh.SendRawPacket(packet, keyNonce)

	//Stampa di tutta la struttura dell'albero per debug
	root.PrintChannelTree()
	return info
}

func testRestoreChannelTree(info ChannelInfo, statePsw string, mainnet bool, keyNonce *KeyNonce) {
	//Si Importano i primi tre layer utilizzando gli stati salvati nel tangle
	root := ImportRootChannelFromTangle(info, statePsw, mainnet)

	//Si prova a creare un nuovo daily channel a partire da un actor esistente ma con data nuova
	root.GetCreateDailyChannelManager(Trucks, "XASD", statePsw, 27, 5, 2021)
	//Si prova a reimportare un daily channel esistente
	dailyCh := root.GetCreateDailyChannelManager(Trucks, "XASD", statePsw, 25, 5, 2021)
	defer dailyCh.Drop()

	//Si riprova ad inviare un altro messaggio di prova
	public, _ := json.Marshal(newMessage("XASD", true))
	private, _ := json.Marshal(newMessage("XASD", false))
	packet := NewRawPacket(public, private)
	defer packet.Drop()
	dailyCh.SendRawPacket(packet, keyNonce)
	root.PrintChannelTree()
}

func main() {
	keyNonce := NewEncryptionKeyNonce("This is a secret key", "This is a secret nonce")
	defer keyNonce.Drop()
	statePsw := "psw"
	mainnet := false
	info := testCreateChannelTree(statePsw, mainnet, keyNonce)
	testRestoreChannelTree(info, statePsw, mainnet, keyNonce)
}
