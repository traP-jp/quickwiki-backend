FROM python:3.10-slim

ENV ARCH amd64
ENV GO_VERSION 1.22.6
# need to change keyword_extractor.go:5 #cgo LDFLAGS: -python{version}

RUN apt update  \
    && apt-get install -y g++ python3-dev wget tar git pkg-config

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

# Build keyword_extractor.cpp
# RUN g++ ./tag/keyword_extractor.cpp -shared -o ./tag/libkeyword_extractor.so -fPIC `pkg-config python3 --cflags` -lpython3.10

# Build Go
RUN go mod tidy && go build -v -o ./main
# WORKDIR /src/tag
# RUN g++ main.cpp -o main.o -L. -lkeyword_extractor `pkg-config python3 --cflags` -lpython$PYTHON_VERSION

# ENV LD_LIBRARY_PATH=/src/tag:$LD_LIBRARY_PATH

ENTRYPOINT ["./main"]
# CMD ["go", "run", "./tag/keyword_extractor.go"]