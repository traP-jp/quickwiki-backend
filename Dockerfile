FROM python:3.10-slim

ENV ARCH amd64
ENV GO_VERSION 1.22.6
# need to change keyword_extractor.go:5 #cgo LDFLAGS: -python{version}

RUN apt update  \
    && apt-get install -y g++ python3-dev wget tar git pkg-config mecab libmecab-dev mecab-ipadic-utf8

# Install Go
RUN cd /tmp  \
    && wget https://golang.org/dl/go$GO_VERSION.linux-$ARCH.tar.gz \
    && tar -C /usr/local -xzf go$GO_VERSION.linux-$ARCH.tar.gz \
    && rm /tmp/go$GO_VERSION.linux-$ARCH.tar.gz

# Set Go environment variables \
ENV PATH="/usr/local/go/bin:${PATH}"


# Install Python dependencies
RUN pip install --upgrade pip
RUN pip install git+https://github.com/boudinfl/pke.git
RUN pip install nltk==3.9.1
RUN pip install -U spacy
RUN python3 -m spacy download en_core_web_sm
RUN python3 -m spacy download ja_core_news_sm
RUN python3 -m nltk.downloader stopwords

COPY ./src /src
WORKDIR /src

ENV CGO_LDFLAGS="-L/usr/lib/x86_64-linux-gnu -lmecab -lstdc++"
ENV CGO_CFLAGS="-I/usr/include"

# Build Go
RUN go mod tidy && go build -v -o ./main

ENTRYPOINT ["./main"]
# CMD ["go", "run", "./tag/keyword_extractor.go"]