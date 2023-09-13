 ## 编译机器
 
 有一台 ARM64 的机器，登录方式如下
 
```shell
ssh root@122.9.143.159

Emqx_123.+
```

进去后，进入 ekuiper 目录, 此目录为 emqx/ekuiper 私有仓库所在地，为不同客户编译不同代码时，自行切换

```shell
cd /root/ekuiper
```

## 中车编译

目前中车 主 branch 为 crrc110, bugfix branch 为 crrc_bugfix。验证阶段的代码均放到了 crrc_bugfix ，待对方验证后合并到 crrc110。
交付方式为 Docker image, Dockerfile 针对中车场景已定制化。 

使用如下命令编译

```shell
docker build  -t ghcr.io/superrxan/ekuiper:crrc-python-arm-Aug-14 -f deploy/docker/Dockerfile-slim-python .
```

编译完后可以打包成 tar, 然后放到某目录让对方下载，对方再本地机器再解压为 docker 镜像

## 启动测试

启动
```shell
docker run -d -p 9081:9081 --name  kuiper-Aug-22 ghcr.io/superrxan/ekuiper:crrc-python-arm-Aug-22
```

进入容器
```shell
 sudo docker exec -it kuiper-Aug-22  /bin/sh
```

查看日志
```shell
docker logs -f kuiper-Aug-22
```

## 交付

### 打包

```shell
docker save -o /tmp/crrc-python-arm-Aug-14.tar  ghcr.io/superrxan/ekuiper:crrc-python-arm-Aug-14
```

### 解包

```shell
docker load -i crrc-python-arm-Aug-14.tar
```

### 下载

```shell
scp root@122.9.143.159:/tmp/crrc-python-arm-Aug-14.tar .
Emqx_123.+
```