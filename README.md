# yatai-arm64

Yatai 1.x 구성요소의 arm64 이미지. 공식 이미지는 amd64 전용이다.

| 이미지 | 태그 |
|---|---|
| `ghcr.io/kurtyoon/yatai-arm64/yatai` | 1.1.13 |
| `ghcr.io/kurtyoon/yatai-arm64/yatai-deployment` | 1.1.21 |
| `ghcr.io/kurtyoon/yatai-arm64/yatai-image-builder` | 1.2.56 |
| `ghcr.io/kurtyoon/yatai-arm64/bento-downloader` | 0.0.5 |
| `ghcr.io/kurtyoon/yatai-arm64/bentoml-proxy` | 0.0.1 |
| `ghcr.io/kurtyoon/yatai-arm64/metrics-transformer` | 0.1.0 |

`images/<이미지>/Dockerfile`이 upstream 소스를 고정 커밋으로 받아 빌드한다. metrics-transformer는 upstream 소스가 공개되어 있지 않아 새로 작성했다. main에 push하면 GitHub Actions가 빌드해 GHCR에 올린다.

upstream(bentoml/Yatai, bentoml/yatai-deployment, bentoml/yatai-image-builder)은 Apache License 2.0이다.
