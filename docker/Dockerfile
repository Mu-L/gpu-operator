# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG GOLANG_VERSION=x.x.x

FROM nvcr.io/nvidia/cuda:12.9.1-base-ubi9 AS builder

RUN dnf install -y wget make git gcc

ARG GOLANG_VERSION=0.0.0
RUN set -eux; \
    \
    arch="$(uname -m)"; \
    case "${arch##*-}" in \
        x86_64 | amd64) ARCH='amd64' ;; \
        ppc64el | ppc64le) ARCH='ppc64le' ;; \
        aarch64 | arm64) ARCH='arm64' ;; \
        *) echo "unsupported architecture" ; exit 1 ;; \
    esac; \
    wget -nv -O - https://storage.googleapis.com/golang/go${GOLANG_VERSION}.linux-${ARCH}.tar.gz \
    | tar -C /usr/local -xz

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor vendor

# Copy the go source
COPY cmd/ cmd/
COPY api/ api/
COPY controllers/ controllers/
COPY internal/ internal/

# Copy Makefile
COPY Makefile Makefile
COPY *.mk .

# Build
ARG VERSION="unknown"
ARG GIT_COMMIT="unknown"
RUN make cmds

FROM nvcr.io/nvidia/cuda:12.9.1-base-ubi9 AS kubectl-downloader

RUN dnf install -y --allowerasing \
        curl \
    && \
        dnf clean all

WORKDIR /workspace

# Install must-gather dependency: `kubectl`
ARG TARGETARCH
RUN OS_ARCH=${TARGETARCH/x86_64/amd64} && OS_ARCH=${OS_ARCH/aarch64/arm64} && \
    curl -LO https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${OS_ARCH}/kubectl && \
    chmod +x ./kubectl

FROM nvcr.io/nvidia/cuda:12.9.1-base-ubi9 AS sample-builder

RUN dnf install -y --allowerasing \
        cmake \
        cuda-cudart-devel-12-9 \
        cuda-nvcc-12-9 \
        g++ \
        curl \
    && \
    dnf clean all

WORKDIR /build

ARG SAMPLE_NAME=vectorAdd

ARG CUDA_SAMPLES_VERSION=v12.9
RUN curl -L https://codeload.github.com/NVIDIA/cuda-samples/tar.gz/refs/tags/${CUDA_SAMPLES_VERSION} | \
    tar -xzvf - --strip-components=1 --wildcards */${SAMPLE_NAME}/* --wildcards */Common/* && \
    cd $(find /build/Samples -iname "${SAMPLE_NAME}") && \
    cmake . && \
    make && \
    cp ${SAMPLE_NAME} /build/${SAMPLE_NAME}

# The C/C++ distroless image is used as a base since the CUDA vectorAdd
# sample application depends on C/C++ libraries.
FROM nvcr.io/nvidia/distroless/cc:v3.1.11-dev

ENV NVIDIA_VISIBLE_DEVICES=void

ARG VERSION="unknown"
ARG GIT_COMMIT="unknown"

LABEL io.k8s.display-name="NVIDIA GPU Operator"
LABEL name="NVIDIA GPU Operator"
LABEL vendor="NVIDIA"
LABEL version="${VERSION}"
LABEL release="N/A"
LABEL summary="Automate the management and monitoring of NVIDIA GPUs."
LABEL description="See summary"
LABEL vsc-ref=${GIT_COMMIT}

WORKDIR /
COPY --from=builder /workspace/gpu-operator /usr/bin/
COPY --from=kubectl-downloader /workspace/kubectl /usr/bin/
COPY --from=builder /workspace/nvidia-validator /usr/bin/
COPY --from=sample-builder /build/vectorAdd /usr/bin/vectorAdd
# TODO: Copy the compat libs from the 'sample-builder' image instead.
# The current 'sample-builder' image does not contain the compat libs in the ARM variant.
# Once new sample images are published that contain the compat libs, we can update the below.
COPY --from=builder /usr/local/cuda/compat /usr/local/cuda/compat

COPY assets /opt/gpu-operator/
COPY manifests /opt/gpu-operator/manifests
COPY validator/manifests /opt/validator/manifests

COPY hack/must-gather.sh /usr/bin/gather

# Add CRD resource into the image for helm upgrades
COPY deployments/gpu-operator/crds/nvidia.com_clusterpolicies.yaml /opt/gpu-operator/nvidia.com_clusterpolicies.yaml
COPY deployments/gpu-operator/crds/nvidia.com_nvidiadrivers.yaml /opt/gpu-operator/nvidia.com_nvidiadrivers.yaml
COPY deployments/gpu-operator/charts/node-feature-discovery/crds/nfd-api-crds.yaml /opt/gpu-operator/nfd-api-crds.yaml

USER 65532:65532

COPY LICENSE /licenses/

ENTRYPOINT ["/usr/bin/gpu-operator"]
