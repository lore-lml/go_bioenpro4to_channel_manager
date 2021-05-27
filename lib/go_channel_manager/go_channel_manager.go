package go_channel_manager

/*
#cgo LDFLAGS: -L./.. -lc_channel_manager_lib
#include "../c_channel_manager.h"
*/
import "C"
import "unsafe"

/*func Hello(str string) {
	cString := C.CString(str)
	C.hello_from_rust(cString)
}*/

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
	key   string
	nonce string
}

type RawPacket struct {
	Public []byte
	Masked []byte
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
	return NewChannelInfo(C.GoString(info.channel_id), C.GoString(info.announce_id))
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

func (root *RootChannel) ChannelInfo() ChannelInfo {
	info := C.root_channel_info(root.root)
	defer C.drop_channel_info(info)
	return NewChannelInfo(C.GoString(info.channel_id), C.GoString(info.announce_id))
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
		kn = keyNonce.toCKeyNonce()
		defer C.drop_key_nonce(kn)
	}

	pack := packet.toCRawPacket()
	defer C.drop_raw_packet(pack)
	var msgId = C.send_raw_packet(ch.channel, pack, kn)
	defer C.drop_str(msgId)
	return C.GoString(msgId)
}

func (ch *DailyChannelManager) ChannelInfo() ChannelInfo {
	info := C.daily_channel_info(ch.channel)
	defer C.drop_channel_info(info)
	return NewChannelInfo(C.GoString(info.channel_id), C.GoString(info.announce_id))
}

func NewChannelInfo(channelId, announceId string) ChannelInfo {
	return ChannelInfo{ChannelId: channelId, AnnounceId: announceId}
}

func NewEncryptionKeyNonce(key, nonce string) *KeyNonce {
	return &KeyNonce{key: key, nonce: nonce}
}

func (keyNonce *KeyNonce) toCKeyNonce() *C.key_nonce_t {
	return C.new_encryption_key_nonce(C.CString(keyNonce.key), C.CString(keyNonce.nonce))
}

func NewRawPacket(pubData, maskData []byte) *RawPacket {
	return &RawPacket{Public: pubData, Masked: maskData}
}

func (packet *RawPacket) toCRawPacket() *C.raw_packet_t {
	cPub, pLen := goByteToCByte(packet.Public)
	cMask, mLen := goByteToCByte(packet.Masked)
	return C.new_raw_packet(cPub, pLen, cMask, mLen)
}

/*func cByteToGoByte(cByte *C.uchar, size C.ulong) []byte{
	return C.GoBytes(unsafe.Pointer(cByte), size)
}*/

func goByteToCByte(bytes []byte) (*C.uchar, C.ulong) {
	return (*C.uchar)(unsafe.Pointer(&bytes[0])), C.ulong(len(bytes))
}
