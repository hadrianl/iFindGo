package ifindgo

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

func TestTHSFUNC(t *testing.T) {
	Initialize("bin")
	login_ec := THS_iFinDLogin(os.Getenv("user"), os.Getenv("pwd"))

	if login_ec != 0 {
		panic(fmt.Errorf("%v errorcode:%v", "login failed!", login_ec))
	}

	ret1 := THS_HistoryQuotes("002632.SZ", "close", "period:D,pricetype:1,rptcategory:0,fqdate:1900-01-01,hb:YSHB,fill:Omit", "2020-09-14", "2020-09-14")
	println(BytesTOString(ret1))
	ret2 := THS_Snapshot("300033.SZ,600000.SH", "preClose;open", "fill:Previous", "2020-09-14 09:30:00", "2020-09-14 10:10:00")
	println(BytesTOString(ret2))
	ret3 := THS_HistoryQuotes("0001.HK,0002.HK,0003.HK,0004.HK,0005.HK,0006.HK", "close", "period:D,pricetype:1,rptcategory:0,fqdate:1900-01-01,hb:YSHB,fill:Previous", "2020-09-14", "2020-09-14")
	println(BytesTOString(ret3))
	ret4 := THS_HighFrequenceSequence("600000.SH", "open;high;low;BBI=BBI_day1:4,BBI_day3:13", "CPS:0,MaxPoints:50000,Fill:Previous,Interval:1", "2020-09-15 09:15:00", "2020-09-15 15:15:00")
	println(BytesTOString(ret4))
	ret5 := THS_EDBQuery("M001889667", "2020-05-19", "2020-05-19")
	println(BytesTOString(ret5))
	ret6 := THS_DateSerial("300033.SZ", "ths_qspj_stock;ths_kpj_stock;ths_zgj_stock;ths_zdj_stock;ths_spj_stock", "100,2020-09-14;100,2020-09-14;100,2020-09-14;100,2020-09-14;100,2020-09-14", "Days:Tradedays,Fill:Previous,Interval:D", "2020-09-01", "2020-09-14")
	println(BytesTOString(ret6))
	ret7 := THS_DateQuery("SSE", "dateType:trade,period:D,dateFormat:0", "2020-07-21", "2020-08-21")
	println(BytesTOString(ret7))
	ret8 := THS_DateOffset("SSE", "dateType:trade,period:W,offset:-10,dateFormat:0", "2020-08-21")
	println(BytesTOString(ret8))
	ret9 := THS_DateCount("SSE", "dateType:trade,period:D,dateFormat:0", "2020-07-21", "2020-08-21")
	println(BytesTOString(ret9))
	ret10 := THS_DataStatistics()
	println(BytesTOString(ret10))
	ret11 := THS_DataPool("index", "2020-09-11;HSI.HK", "date:Y,thscode:Y,security_name:Y,weight:Y")
	println(BytesTOString(ret11))
	ret12 := THS_RealtimeQuotes("601012.SH", "tradeTime;latest;volume", "")
	println(BytesTOString(ret12))
	<-time.After(time.Second * 10)

	logout_ec := THS_iFinDLogout()
	if logout_ec != 0 {
		panic(fmt.Errorf("%v errorcode:%v", "logout failed!", logout_ec))
	}
	println("end")
}

func TestEscape(t *testing.T) {
	Initialize("bin")
	login_ec := THS_iFinDLogin(os.Getenv("user"), os.Getenv("pwd"))

	if login_ec != 0 {
		panic(fmt.Errorf("%v errorcode:%v", "login failed!", login_ec))
	}

	var ptr uintptr
	{
		buf1 := THS_HistoryQuotes("002632.SZ", "close", "period:D,pricetype:1,rptcategory:0,fqdate:1900-01-01,hb:YSHB,fill:Omit", "2020-09-14", "2020-09-14")
		println(BytesTOString(buf1))
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&buf1))
		ptr = sliceHeader.Data
	}

	runtime.GC()

	// should panic
	buf2 := makeByteSlice(ptr)
	println(BytesTOString(buf2))

}
