package ifindgo

/*
#include <string.h>
*/
import "C"
import (
	"errors"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/axgle/mahonia"
)

var h syscall.Handle
var procAddrs map[string]uintptr
var decoder mahonia.Decoder
var mutex sync.Mutex

type errorcode = int32
type queryid = int32

type FT_CALLBACKResultsFunc func(User string, iQueryID int, Result string, ErrorCode, Reserved int) int

func Initialize(iFinDDir string) {
	var err error
	log.SetPrefix("iFinDGO: ")
	bit := 32 << (^uint(0) >> 63)
	// iFinDDir := os.Getenv("iFinDDir")

	log.Println("iFinDDir: ", iFinDDir)
	// find library
	var libPath string
	switch runtime.GOOS {
	case "windows":
		if bit == 32 {
			libPath = "x86\\ShellExport.dll"
		} else {
			libPath = "x64\\ShellExport.dll"
		}

		decoder = mahonia.NewDecoder("gbk")
	case "linux":
		if bit == 32 {
			libPath = "x86/ShellExport.so"
		} else {
			libPath = "x64/ShellExport.so"
		}

		decoder = mahonia.NewDecoder("utf-8")
	default:
		log.Panicln(errors.New("operating systime:" + runtime.GOOS + " not supported"))
	}

	h, err = syscall.LoadLibrary(filepath.Join(iFinDDir, libPath))
	if err != nil {
		log.Panicln(err)
	}

	procnames := []string{
		"THS_BasicData",
		"THS_DataPool",
		"THS_DateSequence",
		"THS_EDBQuery",
		"THS_HighFrequenceSequence",
		"THS_HistoryQuotes",
		"THS_RealtimeQuotes",
		"THS_iFinDLogin",
		"THS_iFinDLogout",
		"THS_AsyBasicData",
		"THS_AsyDataPool",
		"THS_AsyDateSequence",
		"THS_AsyEDBQuery",
		"THS_AsyHighFrequenceSequence",
		"THS_AsyHistoryQuotes",
		"THS_AsyRealtimeQuotes",
		"THS_DeleteBuffer",
		"THS_DataStatistics",
		"THS_GetErrorInfo",
		"THS_DateQuery",
		"THS_DateOffset",
		"THS_DateCount",
		"THS_Snapshot",
		"THS_AsySnapshot",
		"THS_iwencai",
		"THS_Asyiwencai",
		"THS_QuotesPushing",
		"THS_DateSerial",
		"THS_AsyDateSerial",
	}

	procAddrs = make(map[string]uintptr, 29)

	for _, n := range procnames {
		proc, err := syscall.GetProcAddress(h, n)
		if err != nil {
			panic(err)
		}

		procAddrs[n] = proc
	}

}

func THS_BasicData(codes, indicators, parameters string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_BasicData"], uintptr(4),
		uintptr(unsafe.Pointer(s2bp(codes))),
		uintptr(unsafe.Pointer(s2bp(indicators))),
		uintptr(unsafe.Pointer(s2bp(parameters))),
		uintptr(unsafe.Pointer(&ret)), 0, 0)

	return makeByteSlice(ret)
}

func THS_DataPool(datapool, inputparam, outputparam string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_DataPool"], uintptr(4),
		uintptr(unsafe.Pointer(s2bp(datapool))),
		uintptr(unsafe.Pointer(s2bp(inputparam))),
		uintptr(unsafe.Pointer(s2bp(outputparam))),
		uintptr(unsafe.Pointer(&ret)), 0, 0)

	return makeByteSlice(ret)
}

func THS_DateSequence(thscode, jsonIndicator, jsonParam, beginTime, endTime string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_DateSequence"], uintptr(6),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		uintptr(unsafe.Pointer(&ret)))

	return makeByteSlice(ret)
}

func THS_EDBQuery(indicator, beginTime, endTime string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_EDBQuery"], uintptr(4),
		uintptr(unsafe.Pointer(s2bp(indicator))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		uintptr(unsafe.Pointer(&ret)), 0, 0)

	return makeByteSlice(ret)
}

func THS_HighFrequenceSequence(thscode, jsonIndicator, jsonParam, beginTime, endTime string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_HighFrequenceSequence"], uintptr(6),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		uintptr(unsafe.Pointer(&ret)))

	return makeByteSlice(ret)
}

func THS_HistoryQuotes(thscode, jsonIndicator, jsonParam, beginDate, endDate string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_HistoryQuotes"], uintptr(6),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(beginDate))),
		uintptr(unsafe.Pointer(s2bp(endDate))),
		uintptr(unsafe.Pointer(&ret)))

	return makeByteSlice(ret)
}

// TODO: window环境下依赖问题未解决
func THS_RealtimeQuotes(thscode, jsonIndicator, jsonParam string) []byte {
	var ret uintptr

	r, _, err := syscall.Syscall6(procAddrs["THS_RealtimeQuotes"], uintptr(4),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(&ret)), 0, 0)

	println(errorcode(r), err)
	return makeByteSlice(ret)
}

func THS_iFinDLogin(userID, password string) errorcode {

	r, _, _ := syscall.Syscall(procAddrs["THS_iFinDLogin"], uintptr(2), uintptr(unsafe.Pointer(s2bp(userID))), uintptr(unsafe.Pointer(s2bp(password))), 0)

	return errorcode(r)
}

func THS_iFinDLogout() errorcode {
	r, _, _ := syscall.Syscall(procAddrs["THS_iFinDLogout"], uintptr(0), 0, 0, 0)

	return errorcode(r)
}

func THS_AsyBasicData(codes, indicators, parameters string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall6(procAddrs["THS_AsyBasicData"], uintptr(6),
		uintptr(unsafe.Pointer(s2bp(codes))),
		uintptr(unsafe.Pointer(s2bp(indicators))),
		uintptr(unsafe.Pointer(s2bp(parameters))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
	)

	return errorcode(r)
}

func THS_AsyDataPool(datapool, jsonParamArr, jsonFuncOptionalValueArr string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall6(procAddrs["THS_AsyDataPool"], uintptr(6),
		uintptr(unsafe.Pointer(s2bp(datapool))),
		uintptr(unsafe.Pointer(s2bp(jsonParamArr))),
		uintptr(unsafe.Pointer(s2bp(jsonFuncOptionalValueArr))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
	)

	return errorcode(r)
}

func THS_AsyDateSequence(thscode, jsonIndicator, jsonParam, beginTime, endTime string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall9(procAddrs["THS_AsyDateSequence"], uintptr(8),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
		0,
	)

	return errorcode(r)
}

func THS_AsyEDBQuery(indicator, beginTime, endTime string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall6(procAddrs["THS_AsyEDBQuery"], uintptr(6),
		uintptr(unsafe.Pointer(s2bp(indicator))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
	)

	return errorcode(r)
}

func THS_AsyHighFrequenceSequence(thscode, jsonIndicator, jsonParam, beginTime, endTime string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall9(procAddrs["THS_AsyHighFrequenceSequence"], uintptr(8),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
		0,
	)

	return errorcode(r)
}

func THS_AsyHistoryQuotes(thscode, jsonIndicator, jsonParam, beginTime, endTime string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall9(procAddrs["THS_AsyHistoryQuotes"], uintptr(8),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
		0,
	)

	return errorcode(r)
}

func THS_AsyRealtimeQuotes(thscode, jsonIndicator, jsonParam string, bOnlyOnce bool, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall9(procAddrs["THS_AsyRealtimeQuotes"], uintptr(8),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(&bOnlyOnce)),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
		0,
		0,
	)

	return errorcode(r)
}

// TODO: this func might not be necessary in golang, m
func THS_DeleteBuffer(bufptr uintptr) {
	syscall.Syscall(procAddrs["THS_DeleteBuffer"], uintptr(1), uintptr(unsafe.Pointer(&bufptr)), 0, 0)
}

func THS_DataStatistics() []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall(procAddrs["THS_DataStatistics"], uintptr(1), uintptr(unsafe.Pointer(&ret)), 0, 0)

	return makeByteSlice(ret)
}

func THS_GetErrorInfo(ec errorcode) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall(procAddrs["THS_GetErrorInfo"], uintptr(2),
		uintptr(unsafe.Pointer(&ec)), uintptr(unsafe.Pointer(&ret)), 0)

	return makeByteSlice(ret)
}

func THS_DateQuery(exchange, params, startDate, endDate string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_DateQuery"], uintptr(5),
		uintptr(unsafe.Pointer(s2bp(exchange))),
		uintptr(unsafe.Pointer(s2bp(params))),
		uintptr(unsafe.Pointer(s2bp(startDate))),
		uintptr(unsafe.Pointer(s2bp(endDate))),
		uintptr(unsafe.Pointer(&ret)), 0)

	return makeByteSlice(ret)
}

func THS_DateOffset(exchange, params, date string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_DateOffset"], uintptr(4),
		uintptr(unsafe.Pointer(s2bp(exchange))),
		uintptr(unsafe.Pointer(s2bp(params))),
		uintptr(unsafe.Pointer(s2bp(date))),
		uintptr(unsafe.Pointer(&ret)), 0, 0)

	return makeByteSlice(ret)
}

func THS_DateCount(exchange, params, startDate, endDate string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_DateCount"], uintptr(5),
		uintptr(unsafe.Pointer(s2bp(exchange))),
		uintptr(unsafe.Pointer(s2bp(params))),
		uintptr(unsafe.Pointer(s2bp(startDate))),
		uintptr(unsafe.Pointer(s2bp(endDate))),
		uintptr(unsafe.Pointer(&ret)), 0)

	return makeByteSlice(ret)
}

func THS_Snapshot(thscode, jsonIndicator, jsonParam, beginTime, endTime string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall6(procAddrs["THS_Snapshot"], uintptr(6),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		uintptr(unsafe.Pointer(&ret)))

	return makeByteSlice(ret)
}

func THS_AsySnapshot(thscode, jsonIndicator, jsonParam, beginTime, endTime string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall9(procAddrs["THS_AsySnapshot"], uintptr(8),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
		0,
	)

	return errorcode(r)
}

func THS_iwencai(doMain, strType string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall(procAddrs["THS_iwencai"], uintptr(3),
		uintptr(unsafe.Pointer(s2bp(doMain))),
		uintptr(unsafe.Pointer(s2bp(strType))),
		uintptr(unsafe.Pointer(&ret)))

	return makeByteSlice(ret)

}

func THS_Asyiwencai(doMain, strType string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall6(procAddrs["THS_Asyiwencai"], uintptr(5),
		uintptr(unsafe.Pointer(s2bp(doMain))),
		uintptr(unsafe.Pointer(s2bp(strType))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
		0,
	)

	return errorcode(r)
}

func THS_QuotesPushing(thscode, indicator string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall6(procAddrs["THS_QuotesPushing"], uintptr(5),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(indicator))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
		0,
	)

	return errorcode(r)
}

func THS_DateSerial(thscode, jsonIndicator, jsonParam, globalParam, beginTime, endTime string) []byte {
	var ret uintptr

	_, _, _ = syscall.Syscall9(procAddrs["THS_DateSerial"], uintptr(7),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(globalParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		uintptr(unsafe.Pointer(&ret)), 0, 0)

	return makeByteSlice(ret)

}

func THS_AsyDateSerial(thscode, jsonIndicator, jsonParam, globalParam, beginTime, endTime string, callback FT_CALLBACKResultsFunc, user string, queryID queryid) errorcode {
	cb := func(pUser uintptr, iQueryID int, pResultsBuff uintptr, iBuffLength, ErrorCode, Reserved int) int {
		User := C.GoString((*C.char)(unsafe.Pointer(pUser)))
		Result := UTF16TOString(pResultsBuff, iBuffLength)
		return callback(User, iQueryID, Result, ErrorCode, Reserved)
	}

	cbptr := syscall.NewCallback(cb)

	r, _, _ := syscall.Syscall9(procAddrs["THS_AsyDateSerial"], uintptr(9),
		uintptr(unsafe.Pointer(s2bp(thscode))),
		uintptr(unsafe.Pointer(s2bp(jsonIndicator))),
		uintptr(unsafe.Pointer(s2bp(jsonParam))),
		uintptr(unsafe.Pointer(s2bp(globalParam))),
		uintptr(unsafe.Pointer(s2bp(beginTime))),
		uintptr(unsafe.Pointer(s2bp(endTime))),
		cbptr,
		uintptr(unsafe.Pointer(s2bp(user))),
		uintptr(unsafe.Pointer(&queryID)),
	)

	return errorcode(r)
}
