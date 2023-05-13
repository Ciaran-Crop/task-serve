package taskOp

/*
	void HelloWorld(const char* s);
*/
import "C"

func DoHelloWorld() {
	C.HelloWorld(C.CString("Hello World!"))
}
