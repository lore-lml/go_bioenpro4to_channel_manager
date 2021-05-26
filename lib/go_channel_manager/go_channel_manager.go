package go_channel_manager

/*
#cgo LDFLAGS: -L./.. -lc_channel_manager_lib
#include "../c_channel_manager.h"
*/
import "C"

func Hello(str string){
	cString := C.CString(str)
	C.hello_from_rust(cString)
}
