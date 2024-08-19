# c-go-pythonのsetup
pythonの -devをインストール
```shell
sudo apt install python3 python3-dev
```
cファイルのライブラリ作成
```shell
gcc -c hello.c -o hello.o
ar rusv libhello.a hello.o
```
Cで実行するときは以下の内容をmain.cとして
```c
#include "hello.h"

int main(int argc, char const *argv[])
{
    pyHello("yata", "mypy", "hello");
    return 0;
}
```
コンパイルコマンド
```shell
gcc main.c -o main.out -L. -lhello `pkg-config python3 --cflags` -lpython3.10
```
これをcgoに直すと
```go
/*
#cgo pkg-config: python-3.12
#cgo LDFLAGS: -L. -lhello -lpython3.12
#include <Python.h>
#include <stdio.h>
#include <string.h>
#include "hello.h"
*/
import "C"
```

cpythonについて
Py_DECREF: NULL許容しない
Py_XDECREF: NULL許容