package ifindgo

/*
#include <string.h>
*/
import "C"
import (
	"fmt"
	"log"
	"reflect"
	"syscall"
	"unsafe"
)

// FIXME: might cause a memory problem because fo GC,
func s2bp(s string) *byte {
	bp, err := syscall.BytePtrFromString(s)
	if err != nil {
		log.Panicln(err)
	}

	return bp
}

func makeByteSlice(arrPtr uintptr) []byte {
	var buf_c []byte
	size := int(C.strlen((*C.char)(unsafe.Pointer(arrPtr))))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&buf_c))
	sliceHeader.Cap = size
	sliceHeader.Len = size
	sliceHeader.Data = arrPtr

	bs := make([]byte, len(buf_c))
	// copy buf from C to GO
	copy(bs, buf_c)
	THS_DeleteBuffer(arrPtr)
	return bs
}

func makeGoString(arrPtr uintptr) (s string) {
	s = C.GoString((*C.char)(unsafe.Pointer(arrPtr)))
	THS_DeleteBuffer(arrPtr)
	return
}

func BytesTOString(origin []byte) string {
	return decoder.ConvertString(string(origin))
}

func UTF16TOString(retsultPtr uintptr, length int) string {
	var ret []uint16
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	sliceHeader.Data = retsultPtr
	sliceHeader.Cap = length
	sliceHeader.Len = length

	result := syscall.UTF16ToString(ret)

	return result
}

func PrintlnCallback(User string, iQueryID int, Result string, ErrorCode, Reserved int) int {
	fmt.Println(User, iQueryID, ErrorCode, Reserved)
	fmt.Println(Result)
	return 0
}
