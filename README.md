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