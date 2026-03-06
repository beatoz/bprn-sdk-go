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

`EventLog`의 리프 인덱스는 **Global Index (gidx)**를 사용 한다.
Global Index는 Header의 리프들(ChannelId, ChaincodeName, TxId)과 추가된 상태에서의 인덱스를 의미한다.

1. ChannelId (gidx 0)
2. ChaincodeName (gidx 1)
3. TxId (gidx 2)
4. Added Leaf 1 (gidx 3)
5. Added Leaf 2 (gidx 4)
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

#### 5. Custom event element Implementation

`merkle.ILeaves` 인터페이스를 구현하여 애플리케이션에 맞는 커스텀 이벤트 요소 타입을 정의할 수 있다.
각 필드가 머클 트리의 개별 리프가 되므로, 필드 단위로 독립적인 포함 증명이 가능하다.

`merkle.ILeaves` 인터페이스를 구현한 타입은 `EventLog.AddLeaves()` 메서드를 통해 이벤트 로그에 추가할 수 있다.

##### ILeaves 인터페이스

```go
type ILeaves interface {
    Leaf(i int) []byte
    Leaves() [][]byte
    LeavesLen() int
}
```

| 메서드 | 설명 | 구현 내용 |
|--------|------|-----------|
| `Leaf(i int) []byte` | i번째 리프 데이터를 반환한다. | 범위를 초과하면 `nil`을 반환하고, 유효하면 `Leaves()[i]`를 반환한다. |
| `Leaves() [][]byte` | 전체 필드를 `[]byte` 슬라이스로 반환한다. | 각 필드를 `[]byte`로 변환하여 슬라이스로 구성한다. 슬라이스 내 순서가 Global Index에 직접 매핑되므로 순서를 고정해야 한다. |
| `LeavesLen() int` | 리프 개수를 반환한다. | `len(Leaves())`를 반환한다. |

##### Example - TransferEventLog

`samples/transer.go` 에 구현된 토큰 전송 이벤트 예시이다. 4개의 필드가 각각 하나의 리프가 된다.
그 외 구현 예시는 `samples/` 디렉토리를 참고한다.

```go
package samples

import (
    "github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
    "github.com/holiman/uint256"
)

type TransferEventLog struct {
    From   []byte
    To     []byte
    Amount *uint256.Int
    Memo   []byte
}

func (s *TransferEventLog) Leaf(i int) []byte {
    if i >= s.LeavesLen() {
        return nil
    }
    return s.Leaves()[i]
}

func (s *TransferEventLog) Leaves() [][]byte {
    return [][]byte{
        s.From,
        s.To,
        s.Amount.Bytes(),
        s.Memo,
    }
}

func (s *TransferEventLog) LeavesLen() int {
    return len(s.Leaves())
}

var _ merkle.ILeaves = (*TransferEventLog)(nil)
```