package go_channel_manager

/*
#cgo LDFLAGS: -L./.. -lc_channel_manager_lib
#include "../c_channel_manager.h"
*/
import "C"
import "unsafe"

func Hello(str string) {
	cString := C.CString(str)
	C.hello_from_rust(cString)
}

type Category int

const (
	Trucks Category = iota
	Scales
	Biocells
)

type RootChannel struct {
	root *C.root_channel_t
}

type DailyChannelManager struct {
	channel *C.daily_channel_t
}

type ChannelInfo struct {
	ChannelId  string
	AnnounceId string
}

type KeyNonce struct {
	keyNonce *C.key_nonce_t
}

type RawPacket struct {
	packet *C.raw_packet_t
}

func NewRootChannel(mainnet bool) *RootChannel {
	var intMainnet int
	if mainnet {
		intMainnet = 1
	} else {
		intMainnet = 0
	}
	root := C.new_root_channel(C.int(intMainnet))
	return &RootChannel{root: root}
}

func ImportRootChannelFromTangle(channelInfo ChannelInfo, channelPsw string, mainnet bool) *RootChannel {
	var intMainnet int
	if mainnet {
		intMainnet = 1
	} else {
		intMainnet = 0
	}

	cChannelId := C.CString(channelInfo.ChannelId)
	cAnnounceId := C.CString(channelInfo.AnnounceId)
	cPsw := C.CString(channelPsw)
	chInfo := C.new_channel_info(cChannelId, cAnnounceId)
	defer C.drop_channel_info(chInfo)
	root := C.import_root_channel_from_tangle(chInfo, cPsw, C.int(intMainnet))
	return &RootChannel{root: root}
}

func (root *RootChannel) Drop() {
	C.drop_root_channel(root.root)
}

func (root *RootChannel) Open(statePsw string) ChannelInfo {
	cStatePsw := C.CString(statePsw)
	var info = C.open_root_channel(root.root, cStatePsw)
	defer C.drop_channel_info(info)
	return ChannelInfo{ChannelId: C.GoString(info.channel_id), AnnounceId: C.GoString(info.announce_id)}
}

func (root *RootChannel) GetCreateDailyChannelManager(category Category, actorId, statePsw string,
	day, month, year int) *DailyChannelManager {
	cActorId := C.CString(actorId)
	cStatePsw := C.CString(statePsw)
	cDay := C.ushort(day)
	cMonth := C.ushort(month)
	cYear := C.ushort(year)
	channel := C.get_create_daily_actor_channel(root.root, C.int(category), cActorId, cStatePsw, cDay, cMonth, cYear)
	return &DailyChannelManager{channel: channel}
}

func (root *RootChannel) PrintChannelTree() {
	C.print_channel_tree(root.root)
}

func (ch *DailyChannelManager) Drop() {
	C.drop_daily_channel_manager(ch.channel)
}

func (ch *DailyChannelManager) SendRawPacket(packet *RawPacket, keyNonce *KeyNonce) string {
	var kn *C.key_nonce_t = nil
	if keyNonce != nil {
		kn = keyNonce.keyNonce
	}
	var msgId = C.send_raw_packet(ch.channel, packet.packet, kn)
	defer C.drop_str(msgId)
	return C.GoString(msgId)
}

func NewChannelInfo(channelId, announceId string) ChannelInfo {
	return ChannelInfo{ChannelId: channelId, AnnounceId: announceId}
}

func NewEncryptionKeyNonce(key, nonce string) *KeyNonce {
	return &KeyNonce{
		keyNonce: C.new_encryption_key_nonce(C.CString(key), C.CString(nonce)),
	}
}

func (keyNonce *KeyNonce) Drop() {
	C.drop_key_nonce(keyNonce.keyNonce)
}

func NewRawPacket(pubData, maskData []byte) *RawPacket {
	p_len := C.ulong(len(pubData))
	m_len := C.ulong(len(maskData))
	c_pub := (*C.uchar)(unsafe.Pointer(&pubData[0]))
	c_mask := (*C.uchar)(unsafe.Pointer(&maskData[0]))

	var packet = C.new_raw_packet(c_pub, p_len, c_mask, m_len)
	return &RawPacket{packet: packet}
}

func (packet *RawPacket) Drop() {
	C.drop_raw_packet(packet.packet)
}
