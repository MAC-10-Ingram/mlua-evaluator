# mLua Evaluator (TDD CLI Tool for MapleWorlds Lua)

메이플스토리 월드(MSW)의 `mLua` 코드를 MSW 외부의 CLI 환경에서 테스트하고 TDD(Test-Driven Development)를 수행할 수 있도록 도와주는 Golang 기반의 평가 툴입니다. 

## 개요
mLua는 MSW 내부에서만 실행 및 컴파일이 가능하여 단위 테스트나 자동화된 CI/CD 적용이 어렵습니다. 이 툴은 `.mlua` 파일에 포함된 커스텀 어노테이션(`@Component`, `@Sync` 등)과 객체지향 키워드(`script`, `property`, `method`)를 표준 Lua 5.3 문법으로 변환(Transpile)한 뒤, 내장된 Lua VM(`gopher-lua`)을 이용해 개발자가 제공한 데이터 셋과 함께 테스트를 수행합니다.

## 설치 및 빌드

```bash
git clone <repository_url>
cd mlua-evaluator
go build -o mlua-evaluator main.go
```

## 사용 방법

테스트하고자 하는 `.mlua` 파일과 테스트 입력, 모킹(Mocking), 예상 결과값이 정의된 `dataset.json` 파일을 파라미터로 제공합니다.

```bash
./mlua-evaluator <mlua_file> <dataset_json>

# 예시
./mlua-evaluator Unit.mlua dataset.json
```

### 1. `dataset.json` 작성 가이드

테스트 데이터 셋은 각 함수(또는 메서드) 단위로 테스트를 정의합니다. 다음의 내용을 포함합니다.
- `name`: 테스트 케이스 이름
- `target_method`: 실행할 메서드 이름 (예: `Unit:TakeDamage`)
- `mocks`: 모킹할 상태. (Lua 코드로 작성하여 글로벌 변수나 다른 스크립트의 메서드를 모킹)
- `inputs`: 대상 메서드에 전달할 파라미터를 `이름:값` 형태의 객체로 정의합니다.
- `asserts`: 실행 완료 후 검증할 속성과 기대 결과값

**참고:** `mlua` 파일에 정의된 메서드 시그니처의 타입(e.g., `number`, `string`)을 기반으로 **자동 타입 검증**이 수행됩니다.

**`dataset.json` 예시:**
```json
{
  "test_cases": [
    {
      "name": "TakeDamage test",
      "target_method": "Unit:TakeDamage",
      "mocks": [
        "Unit.hp = 1000",
        "_UserService = { GetUser = function() return { name = 'dummy' } end }"
      ],
      "inputs": {
        "amount": 200
      },
      "asserts": [
        {
          "actual": "Unit.hp",
          "expected": 800
        }
      ]
    }
  ]
}
```

### 2. 작동 원리
1. **Parser**: `mlua` 파일에서 메서드 시그니처(파라미터 타입 등)를 추출하고, 데코레이터(`@Sync` 등)를 제거한 뒤 표준 Lua 테이블 및 함수 형태로 트랜스파일합니다.
2. **Runner**: Lua VM을 생성한 후 `mLua` 변환 코드를 적재합니다. 그 다음 `dataset.json`에 정의된 `mocks`를 주입하고, `inputs`의 타입을 검증한 뒤 함수를 실행하여 `asserts`를 검증합니다.

## 디렉토리 구조
- `main.go`: CLI 진입점
- `parser/`: mLua 코드를 Lua로 변환하는 정규식 기반 파서
- `runner/`: Gopher-lua를 활용하여 변환된 스크립트를 실행하고 결과값을 Assert하는 러너
- `mlua-spec.md`: 분석된 mLua 문법 사양 정리
