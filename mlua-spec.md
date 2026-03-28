# mLua 문법 사양 (MapleWorlds Lua)

mLua는 메이플스토리 월드(MSW) 플랫폼에서 게임 개발을 위해 사용되는 Lua 기반의 확장 스크립트 언어입니다. 기본적으로 Lua 5.3 문법을 따르며, 게임 엔진 개발을 위한 객체 지향 문법과 네트워크 동기화를 위한 어노테이션(데코레이터)이 추가되었습니다.

## 1. 스크립트 정의 (Script & Component)

mLua 스크립트는 최상단에 스크립트의 성격을 나타내는 어노테이션과 함께 정의됩니다.

```mlua
@Component
script Unit extends Component
```
- `@Component` : 엔티티에 부착할 수 있는 컴포넌트 스크립트
- `@Logic` : 싱글톤 형태의 로직 스크립트
- `script [ClassName] extends [BaseClass]` 형태로 클래스를 상속받아 정의합니다.

## 2. 프로퍼티 (Property)

클래스 멤버 변수(프로퍼티)는 `property` 키워드와 함께 타입, 식별자, 초기값을 명시합니다. 네트워크 동기화가 필요한 경우 어노테이션을 부착합니다.

```mlua
@Sync
property number hp = 500

@TargetUserSync
property Vector3 originPos = Vector3(0,0,0)

property string attackAnimationRUID = ""
property table spawnSkillData = {}
property Entity activeSkill = nil
```
- **주요 어노테이션**:
  - `@Sync` : 서버에서 변경 시 클라이언트로 자동 동기화되는 프로퍼티
  - `@TargetUserSync` : 특정 유저에게만 동기화되는 프로퍼티
- **주요 데이터 타입**: `number`, `string`, `boolean`, `integer`, `table`, `Entity`, `Vector2`, `Vector3`, `void` 등

## 3. 메서드 (Method)

클래스 함수(메서드)는 `method` 키워드와 함께 반환 타입, 식별자, 매개변수를 명시합니다. 실행 공간을 제어하기 위한 어노테이션을 사용할 수 있습니다.

```mlua
@ExecSpace("Server")
method void Spawn(BuffedUnitModel unitData, table skillData)
    self.baseData = unitData
    self.id = unitData.id
    self:InitData()
end

@ExecSpace("ClientOnly")
method void OnUpdate(number delta)
    if self.isSpawnEnd == false then
        return
    end
end
```
- **실행 공간 어노테이션**:
  - `@ExecSpace("Server")` : 클라이언트에서 호출 시 서버로 RPC 요청, 서버에서 호출 시 서버 실행
  - `@ExecSpace("Client")` : 서버에서 호출 시 클라이언트로 RPC 요청, 클라이언트에서 호출 시 클라이언트 실행
  - `@ExecSpace("ServerOnly")` : 오직 서버에서만 실행 가능
  - `@ExecSpace("ClientOnly")` : 오직 클라이언트에서만 실행 가능
- **이벤트 관련 어노테이션**:
  - `@EventSender("Self")` : 해당 엔티티 내의 이벤트를 수신

## 4. 제어문 및 루아 표준 문법

기본적인 조건문(`if`), 반복문(`for`, `while`), 테이블 조작 및 함수 호출은 Lua 5.3 표준을 그대로 따릅니다.
```mlua
if self.curState == "stand" and self.curTarget == nil then
    self:FindTarget()
elseif self.curState == "move" then
    self:MoveToTarget(delta)
end
```
