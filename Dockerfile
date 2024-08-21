FROM python:3.11.9

ENV ARCH amd64
ENV GO_VERSION 1.22.6
# need to change keyword_extractor.go:5 #cgo LDFLAGS: -python{version}
ENV PYTHON_VERSION 3.11

# Install Go
RUN set -x  \
    && cd /tmp  \
    && wget https://golang.org/dl/go$GO_VERSION.linux-$ARCH.tar.gz \
    && tar -C /usr/local -xzf go$GO_VERSION.linux-$ARCH.tar.gz \
    && rm /tmp/go$GO_VERSION.linux-$ARCH.tar.gz

RUN set -x \
    && apt update && apt-get install -y \
    g++ \
    python$PYTHON_VERSION-dev

# Set Go environment variables \
ENV PATH="/usr/local/go/bin:${PATH}"

# Install Python dependencies
RUN pip install --upgrade pip
RUN pip install git+https://github.com/boudinfl/pke.git
RUN pip install -U nltk
RUN python3 -m spacy download en_core_web_sm
RUN python3 -m spacy download ja_core_news_sm

COPY ./src /src
WORKDIR /src

# Build keyword_extractor.cpp
RUN g++ ./tag/keyword_extractor.cpp -shared -o ./tag/libkeyword_extractor.so -fPIC `pkg-config python3 --cflags` -lpython$PYTHON_VERSION

# Build Go
RUN go mod tidy && go build -o ./main
# WORKDIR /src/tag
# RUN g++ main.cpp -o main.o -L. -lkeyword_extractor `pkg-config python3 --cflags` -lpython$PYTHON_VERSION

ENV LD_LIBRARY_PATH=/src/tag:$LD_LIBRARY_PATH

# ENTRYPOINT ./main
ENTRYPOINT ["./main"]