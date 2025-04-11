package main

/*
#include <stdio.h>
void print(char *s) {
	printf(s);
}
*/
import "C"

func main() {

	C.println(C.CString("hello world"))
}
