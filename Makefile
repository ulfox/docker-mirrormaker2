# Makefile

.DEFAULT_GOAL := build-mirrormaker2

WORKDIR="${PWD}"
# Use for go build binar
SRVER?=

include versions

logdir:
	@if [ ! -e .build ]; then mkdir .build; fi

build-kafka: logdir
	@cd kafka \
		&& docker build \
			--build-arg KAFKA_VERSION="${KAFKA_VERSION}" \
			--build-arg RELEASE_DATE="${RELEASE_DATE}" \
			-f Dockerfile.kafka . \
			-t local/kafka:latest  2>&1 | tee -a ${WORKDIR}/.build/kafka.log

build-mirrormaker2-init: logdir
	@cd mirrormaker2 \
		&& docker build -f Dockerfile.base . \
		-t local/mm2-build:latest 2>&1 | tee -a ${WORKDIR}/.build/mirrormaker2-init.log
 
build-mirrormaker2: logdir build-kafka build-mirrormaker2-init
	@cd mirrormaker2 \
		&& docker build \
			--build-arg KAFKA_VERSION="${KAFKA_VERSION}" \
			--build-arg JAVA_VERSION="${JAVA_VERSION}" \
			--build-arg RELEASE_DATE="${RELEASE_DATE}" \
			-f Dockerfile.mirrormaker2 . \
			-t local/mirrormaker2:latest  2>&1 | tee -a ${WORKDIR}/.build/mirrormaker2.log 

