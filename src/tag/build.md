keyword_extracorから共有ライブラリを生成
```shell
g++ keyword_extractor.cpp -shared -o libkeyword_extractor.so -fPIC `pkg-config python3 --cflags` -lpython3.10
```
ライブライパス
```shell
export LD_LIBRARY_PATH=.:$LD_LIBRARY_PATH
```
C-Python 動作確認用
```shell
g++ main.cpp -o main.o -L. -lkeyword_extractor `pkg-config python3 --cflags` -lpython3.10
```

python: spacy model, pkeが必要 

python3-dev をインストール

ライブラリパスとpythonバージョンは適宜調整

cannot find module エラーはpython側の問題の可能性