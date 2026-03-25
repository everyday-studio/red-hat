# GodotSteam

GodotSteam은 **에디터 교체 방식**으로 사용합니다.
기존 Godot 에디터 대신 GodotSteam이 내장된 커스텀 에디터를 설치하여 사용합니다.
`GodotSteam.app` 안에 Steamworks SDK(`libsteam_api.dylib`)가 이미 포함되어 있으므로
별도 SDK 다운로드는 **불필요**합니다.

## 다운로드

[GodotSteam Releases](https://github.com/GodotSteam/GodotSteam/releases) 최신 릴리즈의 Assets에서:

| 플랫폼 | 다운로드할 파일 |
|--------|----------------|
| **macOS** | `macos-g461-s163-gs4171.tar.xz` |
| Windows | `windows-g461-s163-gs4171.zip` (동일 릴리즈) |

> **`-editor` 파일과 `-templates` 파일은 다운로드하지 마세요.**
> `-editor` 파일은 Android 빌드용 에디터입니다.
> `-templates` 파일은 게임 배포용 익스포트 템플릿으로, 개발 단계에서는 불필요합니다.
>
> 파일명 읽는 법: `macos` (플랫폼) · `g461` (Godot 4.6.1) · `s163` (Steamworks SDK 1.63) · `gs4171` (GodotSteam 4.17.1)

## 설치

1. 다운로드한 파일의 압축을 해제합니다.
2. 나온 `GodotSteam.app`을 `/Applications/` 폴더로 이동합니다.
3. 이 앱을 **기존 Godot 에디터 대신** 사용합니다.
4. `red-hat/apps/client/project.godot`을 이 앱으로 열면 됩니다.

## steam_appid.txt

`apps/client/` 폴더에 `steam_appid.txt` 파일이 있어야 Steam 기능이 활성화됩니다.
이 파일은 `.gitignore`에 등록되어 있어 커밋되지 않습니다.

| 상황 | AppID |
|------|-------|
| 개발 / 로컬 테스트 | `480` (Valve 테스트 앱 Spacewar) |
| Steam 오버레이·도전과제 등 Steam 특화 기능 테스트 | `4525130` (실제 AppID) |

현재 파일에는 `480`이 설정되어 있습니다.
