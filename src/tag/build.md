keyword_extracorから共有ライブラリを生成
```shell
g++ keyword_extractor.cpp -shared -o libkeyword_extractor.so -fPIC `pkg-config python3 --cflags` -lpython3.10
```
ライブライパス
```shell
export LD_LIBRARY_PATH=.:$LD_LIBRARY_PATH
```