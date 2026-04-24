FROM golang:1.22-bookworm

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000

ENV DEBIAN_FRONTEND=noninteractive
ENV GOPATH=/home/${USERNAME}/go
ENV PATH=/home/${USERNAME}/go/bin:/usr/local/go/bin:${PATH}
ENV PIP_DISABLE_PIP_VERSION_CHECK=1

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash-completion \
    build-essential \
    ca-certificates \
    curl \
    git \
    jq \
    less \
    make \
    openssh-client \
    python3 \
    python3-build \
    python3-pip \
    python3-setuptools \
    python3-venv \
    rpm \
    sudo \
    unzip \
    xz-utils \
    zip \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid ${USER_GID} ${USERNAME} \
    && useradd --uid ${USER_UID} --gid ${USER_GID} -m -s /bin/bash ${USERNAME} \
    && echo "${USERNAME} ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME} \
    && mkdir -p /workspaces \
    && mkdir -p /home/${USERNAME}/go/bin \
    && mkdir -p /home/${USERNAME}/go/pkg/mod \
    && mkdir -p /home/${USERNAME}/.cache/go-build \
    && mkdir -p /home/${USERNAME}/.cache/pip \
    && chown -R ${USERNAME}:${USERNAME} /workspaces /home/${USERNAME}

USER ${USERNAME}

RUN go install github.com/goreleaser/nfpm/v2/cmd/nfpm@v2.40.0 \
    && go install github.com/goreleaser/goreleaser/v2@v2.1.0
