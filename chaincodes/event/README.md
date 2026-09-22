## Event Structure

### Abstract

Hyperledger Fabric (HLF) 은 트랜잭션에 오직 하나의 이벤트만이 포함될 수 있도록 구현되어,
다수의 이벤트 발생이 가능한 다른 체인에 비하여 이벤트를 이용한 작업에 많은 한계를 갖는다.
이에 자체적인 이벤트 구조를 정의하고 이를 하나의 바이트 스트림으로 인코딩하여 HLF 이벤트로 기록하고,
이 값을 다시 디코딩하여 정의된 이벤트 구조를 활용한 다양한 기능 구현이 가능하도록 한다.

본 문서는 다음 두 제안을 구현한 이벤트 구조 및 사용법을 설명한다.

* [BTIP-16](https://github.com/beatoz/docs/blob/main/BTIPS/btip-16.md) — `EventLog` 컨테이너 구조, `selector`, Global Index(`gidx`) 매핑
* [BTIP-48](https://github.com/beatoz/docs/blob/main/BTIPS/btip-48.md) — 머클 트리 기술 사양. BTIP-16 의 머클 트리 절을 **대체**한다

> **호환성 경고**
> BTIP-48 적용으로 리프·내부 노드 해시에 도메인 prefix 가 도입되고 패딩·`nil` 리프 처리가 바뀌었다.
> **BTIP-48 이전에 산출된 모든 루트 해시는 이 구현으로 재현되지 않는다.** 버전 협상 없이 일괄 전환하는 것을 전제로 한다.

### Merkle Tree (BTIP-48)

머클 트리는 배열 기반 완전 이진 트리다. 리프 개수는 2의 거듭제곱으로 패딩되며, 노드 해시는 세 도메인으로 분리된다.

| 도메인 | 정의 | 의미 |
|--------|------|------|
| leaf | `LeafHash(data) = sha256(0x00 \|\| data)` | 리프 |
| inner | `InnerHash(l, r) = sha256(0x01 \|\| l \|\| r)` | 내부 노드 |
| empty | `EmptyHash() = sha256(nil)` | 원본 리프가 **없는** 패딩 자리 |
| null | `NullHash() = LeafHash(nil)` | 리프 자리는 **있고** 값이 `nil` 인 리프 |

세 도메인의 해시 입력(길이 0 / `0x00` 시작 / `0x01` 시작 65바이트)이 서로소이므로, SHA-256 의 충돌·제2역상 저항성을 가정할 때
내부 노드를 리프로 제시하거나 서로 다른 깊이의 노드를 바꿔치기하는 것이 계산적으로 불가능하다.

* 리프 1개 트리의 루트는 그 리프의 `LeafHash` 이고, 리프 0개 트리의 루트는 `EmptyHash()` 다.
* `nil` 리프와 패딩 자리는 서로 다른 해시를 가지므로, **패딩 자리에 대한 증명은 만들 수 없다.**
* 지원 최대 깊이는 `merkle.MaxMerkleDepth`(30) 이며, 트리는 최대 `2^30` 개의 리프를 담는다.

### Usage

#### 1. Create Event Log

함수형 옵션으로 이벤트 로그를 생성한다.

```go
log := event.NewEventLog(
    event.WithChannelId("channel-id"),
    event.WithChaincodeId("chaincode-name"),
    event.WithTxId(txId),          // []byte
    event.WithSelector(selector),  // 선택. SetElems 가 덮어쓴다
)
```

#### 2. Set Event Elements

`SetElems` 로 `types.IEventElems` 구현체를 넣는다. `Header.Selector` 가 함께 설정되고, 캐시된 머클 트리는 무효화된다.

```go
log.SetElems(&samples.TransferEventElems{
    From:   from,
    To:     to,
    Amount: amount,
    Memo:   memo,
})
```

> `Header` 와 `Elems` 필드를 **직접 수정하면 캐시된 트리가 무효화되지 않는다.**
> 값 변경은 `SetElems` / `Reset` / `UnmarshalDER` 경로로만 수행해야 한다.

#### 3. Marshal and Unmarshal

`MarshalDER` / `UnmarshalDER` 로 DER(ASN.1) 인코딩한다. `onlyElems` 를 주면 Header 를 제외하고 `Elems` 만 처리한다.

```go
data, err := log.MarshalDER()

newLog := &event.EventLog{}
err = newLog.UnmarshalDER(data)

// Elems 만
data, err = log.MarshalDER(true)
err = newLog.UnmarshalDER(data, true)
```

#### 4. Merkle Root and Proof

`EventLog` 의 리프 인덱스는 **Global Index (gidx)** 를 사용한다.

| gidx | 리프 |
|------|------|
| 0 | ChannelId |
| 1 | ChaincodeId |
| 2 | TxId |
| 3 | Selector |
| 4 ~ | `SetElems` 로 추가된 이벤트 요소 |

```go
root := log.Root()

leaf, siblings, err := log.Proof(0)
if err != nil {
    // gidx 가 [0, LeavesLen()) 밖이면 오류
}
```

증명은 `(index, leaf, siblings)` 세 값이다. `Proof` 는 이 중 **`index` 를 돌려주지 않는다** — 호출자가 방금 건넨 인자이기 때문이다. 검증 시 다시 필요하므로 호출자가 들고 있어야 하고, 증명을 전송한다면 세 값을 함께 보내야 한다.

`leaf` 는 리프 해시가 아니라 **원본 데이터**다. 검증기가 항상 `LeafHash` 를 적용하기 때문이다.

**증명 검증 — 두 경로를 구분한다.**

```go
// (1) 자기 일관성 확인: 이 EventLog 자신의 루트로 검증한다.
err := log.VerifyProof(0, leaf, siblings)

// (2) 신뢰된 외부 루트로 검증한다. 실제 증명 검증은 이 경로를 쓴다.
err = merkle.VerifyProof(0, leaf, siblings, trustedRoot) // 예: 서명된 block_event_root 에서 유도한 루트
```

`merkle.VerifyProof` 는 접기 전에 다음을 모두 검사한다.

* `len(siblings) <= merkle.MaxMerkleDepth`
* `root` 와 모든 sibling 이 정확히 32바이트
* 증명 깊이로 복원한 `n = 1 << len(siblings)` 에 대해 `0 <= index < n`

마지막 검사가 없으면 `index` 의 상위 비트가 무시되어 `index` 와 `index + n` 이 같은 증명으로 통과한다. `index` 가 구조체에 묶여 오지 않더라도 이 검사는 검증기 입력에 대해 그대로 적용되므로 방어는 유지된다.

#### 5. 슬라이스 소유권

`merkle.MerkleTree` 는 넘겨받은 리프를 **복사하지 않고 그대로 들고 있는다**. 생성 후 그 슬라이스를 수정하면 루트는 그대로지만 `Proof` 가 트리가 커밋한 적 없는 리프를 돌려주므로 검증이 실패한다.

반면 `Root()` 와 `Proof()` 가 돌려주는 값은 **복사본**이므로 호출자가 마음대로 써도 된다. 패딩 슬롯과 `nil` 리프는 패키지 전역 테이블을 가리키므로 이 복사가 반드시 필요하다.

자세한 내용은 [docs/merkle-slice-ownership.md](./docs/merkle-slice-ownership.md) 를 참고한다.

#### 6. Custom Event Element Implementation

`types.IEventElems`(= `types.ILeaves` + `Selector()`) 를 구현하면 애플리케이션 고유의 이벤트 요소 타입을 정의할 수 있다.
각 필드가 머클 트리의 개별 리프가 되므로 필드 단위 포함 증명이 가능하다.

```go
type ILeaves interface {
    Leaf(i int) []byte
    Leaves() [][]byte
    LeavesLen() int
}

type IEventElems interface {
    ILeaves
    Selector() []byte
}
```

| 메서드 | 구현 내용 |
|--------|-----------|
| `Leaf(i)` | 범위를 벗어나면 `nil`, 유효하면 `Leaves()[i]` |
| `Leaves()` | 각 필드를 `[]byte` 로 변환한 슬라이스. **순서가 곧 gidx 이므로 고정해야 한다** |
| `LeavesLen()` | `len(Leaves())` |
| `Selector()` | 이벤트 시그니처 해시. `SetElems` 가 `Header.Selector` 에 넣는다 |

구현 예시는 `samples/` 디렉토리를 참고한다 (`samples/transfer_evtelems.go`, `samples/post_message.go`).

### API 마이그레이션 (BTIP-48)

| 이전 | 현재 |
|------|------|
| `merkle.WithHashedLeaves(leaves)` | **삭제.** 모든 리프는 `LeafHash` 를 거친다. 32바이트 해시를 리프로 쓰더라도 `WithRawLeaves` 로 넘긴다 |
| `VerifyProof(idx, data, siblings, root, preHashed...)` | `merkle.VerifyProof(index int, leaf []byte, siblings [][]byte, root []byte) error` — `preHashed` 없음 |
| `tree.Proof(i) ([]byte, [][]byte, error)` — 첫 값은 리프 **해시** | `tree.Proof(i) ([]byte, [][]byte, error)` — 첫 값은 리프 **원본 데이터** |
| `log.Proof(gidx) ([]byte, [][]byte, error)` | 동일. 단 첫 값이 리프 **원본 데이터** |
| `log.VerifyProof(gidx, siblings)` | `log.VerifyProof(gidx int, leaf []byte, siblings [][]byte) error` — `leaf` 인자 추가 |
| 하위 트리 루트를 상위 트리에 그대로 배치 | 상위 트리의 **원본 리프 데이터**이므로 `LeafHash` 를 다시 거친다 (트리 계층 간 도메인 바인딩) |
