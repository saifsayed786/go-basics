package main

import (
	"fmt"
	"unsafe"
)

type MyStruct struct {
	Field1 int32
	Field2 string
	Field3 bool
}

func structToBytes(data MyStruct) []byte {
	var byteSlice []byte

	// Convert the int field to bytes
	field1Bytes := (*[unsafe.Sizeof(&data.Field1)]byte)(unsafe.Pointer(&data.Field1))[:]
	byteSlice = append(byteSlice, field1Bytes...)

	// Convert the string field to bytes
	field2Bytes := []byte(data.Field2)
	byteSlice = append(byteSlice, field2Bytes...)

	// Convert the bool field to a single byte
	var field3Byte byte
	if data.Field3 {
		field3Byte = 1
	} else {
		field3Byte = 0
	}
	byteSlice = append(byteSlice, field3Byte)

	return byteSlice
}

func main() {
	MemorySize()
	data := MyStruct{
		Field1: 42,
		Field2: "Hello, World!",
		Field3: true,
	}

	bytes := structToBytes(data)
	fmt.Println(bytes) // Output: [42 0 0 0 72 101 108 108 111 44 32 87 111 114 108 100 33 1]
	restoredData := bytesToStruct(bytes)
	fmt.Println(restoredData)
}

func MemorySize() {
	var i int

	var i16 int16
	var i32 int32
	var i64 int64

	fmt.Printf("i Type:%T Size:%d\n", i, unsafe.Sizeof(i))
	fmt.Printf("i16 Type:%T Size:%d\n", i16, unsafe.Sizeof(i16))
	fmt.Printf("i32 Type:%T Size:%d\n", i32, unsafe.Sizeof(i32))
	fmt.Printf("i64 Type:%T Size:%d\n", i64, unsafe.Sizeof(i64))
}

func bytesToStruct(bytes []byte) MyStruct {
	var data MyStruct

	// Convert the first 4 bytes to int32
	data.Field1 = *(*int32)(unsafe.Pointer(&bytes[0]))

	// Convert the remaining bytes to string
	data.Field2 = string(bytes[4 : len(bytes)-1])

	// Convert the last byte to bool
	data.Field3 = bytes[len(bytes)-1] != 0
	fmt.Println(data.Field1)
	return data
}
