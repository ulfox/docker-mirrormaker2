# Docker Mirrormaker2

Bootstrap MM2 docker containers.

##  Build

To create a Mirrormaker2 docker container first update the `versions` file
with the desried `KAFKA & JAVA` versions

### Versions

Update file `versions` with the desired versions

```bash
KAFKA_VERSION="2.8.1"
JAVA_VERSION="11"
RELEASE_DATE="2021-09-19"
```

### Build Mirrormaker2 Docker Container

    make

#### Build logs

You can find build logs for debug under .build

```bash
$> ls .build/
kafka.log  mirrormaker2-init.log  mirrormaker2.log
```

## Usage

We can configure and run MM2 with multiple ways

### Using a source mm2 properties file


    version: '3'
    services:
      mm2:
        image: local/mirrormaker2:latest
        container_name: mirrormaker2
        volumes:
            - /path/to/my/kafka-mm2.properties:/opt/mm2/kafka-mm2.properties:ro


### Using env variables

    version: '3'
    services:
      mm2:
        image: local/mirrormaker2:latest
        container_name: mirrormaker2
        env:
            KMM2_SOME_KEY_1: SOME_VALUE_1
            ...
            KMM2_SOME_KEY_N: SOME_VALUE_N

**Example**

    version: '3'
    services:
      mm2:
        image: local/mirrormaker2:latest
        container_name: mirrormaker2
        environment:
          KMM2_CLUSTERS: source, target
          KMM2_BOOTSTRAP_SERVERS: PLAINTEXT://localhost:9092
          ...

### Mixed Env and Source file

    version: '3'
    services:
      mm2:
        image: local/mirrormaker2:latest
        container_name: mirrormaker2
        env:
            KMM2_SOME_KEY_1: SOME_VALUE_1
            ...
            KMM2_SOME_KEY_N: SOME_VALUE_N
        volumes:
            - /path/to/my/kafka-mm2.properties:/opt/mm2/kafka-mm2.properties:ro


**Example**

    version: '3'
    services:
      mm2:
        image: local/mirrormaker2:latest
        container_name: mirrormaker2
        environment:
          KMM2_CLUSTERS: source, target
          KMM2_BOOTSTRAP_SERVERS: PLAINTEXT://localhost:9092
        volumes:
            - ./mirrormaker2/example/kafka-mm2.properties:/opt/mm2/kafka-mm2.properties:ro