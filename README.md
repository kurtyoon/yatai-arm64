# yatai-arm64

Yatai 1.x 구성요소를 arm64로 빌드하는 저장소다. Yatai 1.x의 공식 이미지는 모두 amd64 전용이고 `bentoml/Yatai` 저장소는 아카이브되었다. arm64 단일 노드 k3s 실습 환경에서 Yatai를 한 번 돌려 보기 위해 만들었다.

## 이미지

모두 `ghcr.io/kurtyoon/yatai-arm64/<이미지>`에 올라간다. 태그는 대응하는 upstream 이미지의 태그를 그대로 쓰고, 같은 이미지에 `<태그>-<커밋 7자리>` 태그가 하나 더 붙는다. 배포 매니페스트는 커밋 태그나 digest로 고정한다.

| 이미지 | 태그 | 대응하는 upstream 이미지 | 만드는 방법 |
|---|---|---|---|
| `yatai` | 1.1.13 | `quay.io/bentoml/yatai:1.1.13` | API 서버를 소스에서 빌드, 대시보드 정적 파일은 upstream 이미지에서 복사 |
| `yatai-deployment` | 1.1.21 | `quay.io/bentoml/yatai-deployment:1.1.21` | 소스에서 빌드 |
| `yatai-image-builder` | 1.2.56 | `quay.io/bentoml/yatai-image-builder:1.2.56` | 소스에서 빌드 |
| `bento-downloader` | 0.0.5 | `quay.io/bentoml/bento-downloader:0.0.5` | aarch64 AWS CLI로 다시 패키징, gcloud 제외 |
| `bentoml-proxy` | 0.0.1 | `quay.io/bentoml/bentoml-proxy:0.0.1` | envoy 1.24.1 aarch_64 바이너리로 다시 패키징 |
| `metrics-transformer` | 0.1.0 | `quay.io/bentoml/yatai-bento-metrics-transformer:0.0.4` | 새로 작성 |

kaniko는 `quay.io/bentoml/kaniko:debug`가 이미 arm64를 포함하므로 빌드하지 않는다.

## 만드는 방법

- 각 `images/<이미지>/Dockerfile`이 upstream 저장소를 고정된 커밋으로 받아 빌드한다. 커밋, 버전, Go 빌드 플래그는 upstream Makefile과 동일하고, `GOARCH`만 빌드 대상 플랫폼에서 가져온다.
- Go는 1.27.1로 빌드한다. upstream은 1.18에서 1.20을 썼다. Go 표준 라이브러리에서 오는 취약점은 새 Go로 빌드하면 빠진다. 의존 모듈의 버전은 upstream 그대로다.
- `metrics-transformer`는 upstream 소스가 공개되어 있지 않아서 새로 작성했다. yatai-deployment가 사이드카로 띄우면서 넘기는 환경 변수(`BENTOML_SERVER_HOST`, `BENTOML_SERVER_PORT`, `PORT`, `OLD_METRICS_PREFIX`, `NEW_METRICS_PREFIX`)를 받는다. `/healthz`에 응답하고, Bento의 `/metrics`를 받아 지표 이름의 접두사를 바꿔 돌려준다.
- `bento-downloader`는 yatai-image-builder가 init 컨테이너로 띄워 Bento와 모델을 받는다. 이 환경은 `gs://` 주소를 쓰지 않아 gcloud를 넣지 않았다.

## 빌드

`.github/workflows/build.yaml`이 GitHub의 arm64 runner에서 빌드한다.

- pull request: 모든 이미지를 빌드만 하고 올리지 않는다.
- `main` push와 수동 실행: 빌드한 뒤 GHCR에 올린다.

## 라이선스

upstream 프로젝트(bentoml/Yatai, bentoml/yatai-deployment, bentoml/yatai-image-builder)는 Apache License 2.0이다. API 서버와 두 운영자 이미지는 그 소스를 수정 없이 빌드한 것이다.
