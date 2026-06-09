# Docker 常用命令速查表

> 项目 bbs-go 使用 `--docker` 运行时部署在 k3s 中，Docker 是底层容器引擎

---

## 一、镜像管理

```bash
# 列出所有镜像
docker images

# 按名称过滤
docker images | grep temp-

# 删除镜像
docker rmi temp-auth:old-tag

# 删除所有未被使用的镜像
docker image prune -a -f

# 构建镜像（多阶段构建）
docker build -f auth/Dockerfile -t temp-auth:latest .

# 查看镜像历史层
docker history temp-auth:latest
```

---

## 二、容器管理

```bash
# 列出运行中的容器
docker ps

# 列出所有容器（含已退出的）
docker ps -a

# 列出所有容器（含已退出的，只显示 name 和 status）
docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Image}}"

# 查看容器日志
docker logs <container-name>

# 查看容器日志（实时跟踪）
docker logs -f <container-name>

# 进入容器内部
docker exec -it <container-name> sh

# 停止容器
docker stop <container-name>

# 删除容器
docker rm <container-name>

# 强制删除运行中的容器
docker rm -f <container-name>

# 重启容器
docker restart <container-name>
```

---

## 三、容器清理

```bash
# 删除所有已退出的容器
docker container prune -f

# 删除所有未使用的资源（容器、网络、镜像、构建缓存）
docker system prune -f

# 更彻底：连未使用的镜像一起删
docker system prune -a -f

# 查看占用磁盘
docker system df
```

> ⚠️ `docker system prune -a` 会删除所有未被任何容器引用的镜像，后续用的时候需要重新 build 或 pull。k3s 管理的容器不要手动 `rm`，通过 `kubectl delete pod` 控制。

---

## 四、导入导出（用于 k3s）

```bash
# 将 Docker 镜像导出并导入 k3s containerd
docker save temp-auth:latest | k3s ctr images import -

# 或者保存为 tar 文件
docker save temp-auth:latest -o temp-auth.tar

# 从 tar 文件导入到 k3s
k3s ctr images import temp-auth.tar

# 清理 tar 文件
rm temp-auth.tar
```

---

## 五、网络

```bash
# 列出网络
docker network ls

# 查看网络详情
docker network inspect bridge

# 查看 k3s 默认网络（docker 运行时）
docker network ls | grep k3s
```

---

## 六、常用组合操作

```bash
# 构建 + 导入 k3s（改代码后一条龙）
docker build -f auth/Dockerfile -t temp-auth:latest . && \
docker save temp-auth:latest | k3s ctr images import - && \
kubectl rollout restart deployment/auth -n temp

# 查看容器资源占用
docker stats --no-stream

# 查看某个容器的详细信息
docker inspect <container-id> | grep -A10 "Mounts"
```
