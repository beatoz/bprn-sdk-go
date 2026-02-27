## Event Structure

### Abstract

Hyperledger Fabric (HLF) 은 트랜잭션에 오직 하나의 이벤트만이 포함될 수 있도록 구현되어,
다수의 이벤트 발생이 가능한 다른 체인에 비하여 이벤트를 이용한 작업에 많은 한계를 갖는다.
이에 자체적인 이벤트 구조를 정의하고 이를 하나의 바이트 스트림으로 인코딩하여 HLF 이벤트로 기록하고,
이 값을 다시 디코딩하여 정의된 이밴트 구조를 활용한 다양한 기능 구현이 가능하도록 한다.

본 문서는 [BTIP-16](https://github.com/beatoz/docs/blob/main/BTIPS/btip-16.md) 제안 내용을 구현한 이벤트 구조 및 사용법을 설명한다.

### Usage

#### 1. Create Event Log

`NewEventLog` 함수를 사용하여 이벤트 로그를 생성한다.

```go
log := event.NewEventLog("channel-id", "chaincode-name", "tx-id")
```

#### 2. Add Event Data

`AddData` 또는 `AddLeaves` 메서드를 사용하여 이벤트 데이터를 추가한다.

*   `AddData`: 단일 바이트 배열을 리프로 추가한다.
*   `AddLeaves`: `types.ILeaves` 인터페이스를 구현한 객체의 모든 리프를 추가한다.

```go
// Add single leaf
log.AddData([]byte("some-data"))

// Add multiple leaves from an object implementing ILeaves
type MyData struct {
    // ...
}
func (d *MyData) Leaves() [][]byte {
    // ...
}

data := &MyData{...}
log.AddLeaves(data)
```

#### 3. Marshal and Unmarshal

`MarshalEventLog` 함수를 사용하여 이벤트 로그를 바이트 스트림으로 인코딩하고, `UnmarshalEventLog` 함수를 사용하여 디코딩한다.

```go
// Marshal
data, err := event.MarshalEventLog(log)
if err != nil {
    // handle error
}

// Unmarshal
newLog, err := event.UnmarshalEventLog(data)
if err != nil {
    // handle error
}
```

#### 4. Merkle Tree

`EventLog`는 내부적으로 머클 트리를 구성하여 루트 해시를 계산하고, 특정 필드에 대한 머클 증명을 생성 및 검증할 수 있다.

**Calculate Root Hash:**

```go
// Get Root Hash
root := log.Root()
```

**Generate Merkle Proof:**

`EventLog`의 리프 인덱스는 **Global Index**를 사용해야 한다.
Global Index는 Header의 리프들(ChannelId, ChaincodeName, TxId)과 추가된 데이터 리프들을 순서대로 포함한다.

1. ChannelId (Index 0)
2. ChaincodeName (Index 1)
3. TxId (Index 2)
4. Added Leaf 1 (Index 3)
5. Added Leaf 2 (Index 4)
...

```go
// Generate Proof for ChannelId (Global Index 0)
leafHash, proof, err := log.Proof(0)
if err != nil {
    // handle error
}
```

**Verify Merkle Proof:**

```go
// Verify Proof
err := log.VerifyProof(0, proof)
if err == nil {
    // Proof is valid
} else {
    // Proof is invalid
}
```
