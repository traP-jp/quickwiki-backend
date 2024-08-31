package main

/*
const char* hello() {
	const char *text = "Hello, World!";
	return text;
}
*/
import "C"
import (
	"log"
)

func main() {
	text := C.hello()
	log.Println(C.GoString(text))
}
