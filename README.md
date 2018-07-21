# gRPC-gacha

gRPCの勉強として、[gRPCを使った簡易的なマイクロサービスを作ってみた](https://qiita.com/kotamat/items/a84301a16fc24a203304)を触ってみたメモ。


## protobufの定義、コンパイル

```proto
syntax = "proto3";

package gacha;

service Gacha {
  rpc Lottery (Request) returns (Response) {}
}

message Card {
  string name = 1;
}

message Request {
  repeated Card cards = 1;
}

message Response {
  Card card = 1;
  int32 ret_code = 2;
}
```

```sh
❯ cd gacha
❯ protoc --go_out=plugins=grpc:./ gacha.proto
```

`gacha`ディレクトリの中に`gacha.pb.go`が生成される。

## 実装

### server

`gacha.proto`の`service`で定義したrpc`Lottery`を実装する必要がある。  
`gacha.pb.go`で定義されるinterfaceは以下の通り。

```go
// gacha.pb.go
type GachaServer interface {
	Lottery(context.Context, *Request) (*Response, error)
}
```

`GachaServer`のinterfaceを実装した`server`を以下のようにgrpcのサーバに登録する。

```go
s := grpc.NewServer()
pb.RegisterGachaServer(s, &server{})
```

あとはlistenで待ち受ける。
サーバ全体の実装は以下の通り。

```go
// server.go
package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/cipepser/gRPC-gacha/gacha"
	"google.golang.org/grpc"
)

type server struct {
}

const (
	port = ":8080"
)

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGachaServer(s, &server{})
	s.Serve(l)
}

func (s *server) Lottery(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	if len(in.Cards) < 1 {
		return nil, errors.New("empty cards")
	}
	rand.Seed(time.Now().UnixNano())
	chosenKey := rand.Intn(len(in.Cards))
	return &pb.Response{Card: in.Cards[chosenKey], RetCode: 1}, nil
}
```

### client

クライアントは、`grpc.Dial`で接続しに行く。
クライアントから`Lottery`を呼び出すのでサーバ側と同じinterfaceになるように実装する。

```go
// client.go
package main

import (
	"context"
	"log"

	pb "github.com/cipepser/gRPC-gacha/gacha"
	"google.golang.org/grpc"
)

const (
	address = "localhost"
	port    = ":8080"
)

type client struct {
}

func main() {
	conn, err := grpc.Dial(address+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGachaClient(conn)

	cards := []*pb.Card{
		&pb.Card{Name: "card1"},
		&pb.Card{Name: "card2"},
		&pb.Card{Name: "card3"},
	}

	resp, err := c.Lottery(context.Background(), &pb.Request{Cards: cards})
	if err != nil {
		log.Fatalf("could not get card %v", err)
	}
	log.Printf("get card: %v", resp.Card.Name)
}
```

## 実行

### server

```sh
❯ go run server/server.go
```

### client

```sh
❯ go run client/client.go
2018/07/21 17:28:37 get card: card1

❯ go run client/client.go
2018/07/21 17:28:40 get card: card3

❯ go run client/client.go
2018/07/21 17:28:42 get card: card1

❯ go run client/client.go
2018/07/21 17:28:45 get card: card3

❯ go run client/client.go
2018/07/21 17:28:47 get card: card1

❯ go run client/client.go
2018/07/21 17:28:49 get card: card3

❯ go run client/client.go
2018/07/21 17:28:50 get card: card3

❯ go run client/client.go
2018/07/21 17:28:52 get card: card2
```

## References
* [gRPCを使った簡易的なマイクロサービスを作ってみた](https://qiita.com/kotamat/items/a84301a16fc24a203304)
* [Protocol Buffer Basics: Goの写経](https://github.com/cipepser/gRPC-sample/tree/master/1)