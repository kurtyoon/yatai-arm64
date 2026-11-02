# bento-downloader 0.0.5 for arm64: bento-downloader/Dockerfile from
# yatai-image-builder v1.2.56, with the aarch64 AWS CLI.
# yatai-image-builder runs it as an init container that fetches the Bento
# (curl, aws s3 cp) and models (huggingface-cli). gcloud is left out: it only
# serves gs:// URLs, which this lab does not use.

FROM ubuntu:22.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
      ca-certificates curl jq unzip xz-utils file python3 python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Upstream installed the latest huggingface_hub. 0.26.5 still ships the
# huggingface-cli command the init script calls.
RUN pip3 install --no-cache-dir "huggingface_hub[cli]==0.26.5"

RUN curl -fsSL "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o /tmp/awscliv2.zip \
    && unzip -q /tmp/awscliv2.zip -d /tmp \
    && /tmp/aws/install \
    && rm -rf /tmp/awscliv2.zip /tmp/aws

# Ahead of /usr/bin in PATH: the init script's curl calls retry refused
# connections (see the file).
COPY --chmod=0755 curl /usr/local/bin/curl
