FROM golang:1.24

WORKDIR /app

RUN apt-get update && apt-get install -y protobuf-compiler

ENV PATH="/go/bin:$PATH"

COPY . .

RUN make install
RUN make generate
RUN make tidy
CMD ["make", "run"]
