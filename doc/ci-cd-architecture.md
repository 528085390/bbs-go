# CI/CD 流水线架构详解

> 项目：bbs-go | 框架：go-zero | 集群：K3s | GitOps：ArgoCD

---

## 一、整体架构图

```
                    ┌─────────────────────────────────────────────┐
                    │              开发者                          │
                    │          git push (改代码)                   │
                    └──────────────────┬──────────────────────────┘
                                       │
                                       ▼
┌──────────────────────────────────────────────────────────────────┐
│                  ① GitHub Actions (CI)                          │
│                                                                  │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────────┐ │
│  │ 测试阶段  │ → │ 构建阶段  │ → │ 推送阶段  │ → │ 配置更新阶段  │ │
│  │ go test  │   │ docker   │   │ push 到  │   │ 改 kustomize │ │
│  │ ./...    │   │ build    │   │ ghcr.io  │   │ 镜像版本号   │ │
│  └──────────┘   └──────────┘   └──────────┘   └──────┬───────┘ │
│                                                       │         │
│                                         git commit "更新镜像标签"│
│                                         push 回 GitHub 仓库     │
└──────────────────────────────────────────────────────┬───────────┘
                                                       │
                                                       ▼
┌──────────────────────────────────────────────────────────────────┐
│                  ② GitHub 仓库 (配置仓库)                        │
│                                                                  │
│  deploy/k8s/overlays/dev/kustomization.yaml                      │
│    镜像版本从 :latest → :sha-abc123                              │
└──────────────────────────────────────────────────────┬───────────┘
                                                       │
                          ArgoCD 每 3 分钟轮询一次 ─────┘
                                                       ▼
┌──────────────────────────────────────────────────────────────────┐
│                  ③ ArgoCD (CD)                                  │
│                                                                  │
│  检测到 Git 仓库的 YAML 文件变化                                   │
│  → 对比当前集群状态 → 执行同步                                    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐     │
│  │ 同步策略：Auto-Sync + Prune + Self-Heal                  │     │
│  │ • Auto-Sync: 自动应用变更，无需手动 kubectl apply        │     │
│  │ • Prune: 删除 Git 中已移除的资源                          │     │
│  │ • Self-Heal: 有人手动改集群，ArgoCD 会改回 Git 定义的版本│     │
│  └─────────────────────────────────────────────────────────┘     │
└──────────────────────────────────────────────────────┬───────────┘
                                                       │
                                                       ▼
┌──────────────────────────────────────────────────────────────────┐
│                  ④ K3s 集群                                      │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐     │
│  │ 滚动更新过程：                                            │     │
│  │ 1. 创建新 RS（新版本），启动新 Pod                        │     │
│  │ 2. 新 Pod 就绪后，停止旧 Pod                              │     │
│  │ 3. 逐步替换，全程服务在线                                  │     │
│  │ 4. 旧 RS 保留 0 个 Pod，用于回滚                          │     │
│  └─────────────────────────────────────────────────────────┘     │
│                                                                  │
│  最终状态：9 个微服务 + 3 个基础设施全部 Running                   │
└──────────────────────────────────────────────────────────────────┘
```

---

## 二、各组件详解

### 1. GitHub Actions（持续集成 CI）

**职责**：代码推送到 GitHub 后，自动执行测试、构建、推送镜像、更新配置。

**工作流文件**：`.github/workflows/ci.yml`

**包含 3 个 Job：**

#### Job 1: test（测试）
```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.25'
- run: go test ./...
```
- 作用：运行 Go 单元测试，保证代码质量
- 失败后：整个流水线停止，不继续构建
- 为什么需要：防止有 bug 的代码部署到生产

#### Job 2: build（构建 + 推送）
```yaml
strategy:
  matrix:
    service: [auth, user, section, post, ...]
```
- 并行构建 9 个服务的 Docker 镜像（Job 2 依赖 Job 1）
- 每个镜像打两个标签：
  - `sha-<git commit hash>`（精确版本）
  - `latest`（最新版本）
- 推送到 ghcr.io（GitHub Container Registry）

**为什么用 matrix？**—— 9 个服务并行构建，总时间 ≈ 单个服务的构建时间（～3分钟），而不是顺序构建的9倍时间。

#### Job 3: update-config（更新配置）
```yaml
kustomize edit set image temp-auth=ghcr.io/.../auth:sha-xxx
git commit -m "chore: update image tags [skip ci]"
git push
```
- 修改 `deploy/k8s/overlays/dev/kustomization.yaml` 中的镜像标签
- 从 `:latest` 改为 `:sha-<commit>`，实现 **不可变部署**（每次部署都是确定性的版本）
- 用 `[skip ci]` 避免触发无限循环（CI push → 触发 CI → push → ...）

---

### 2. GitHub Container Registry（镜像仓库）

**地址**：`ghcr.io/528085390/bbs-go/<service>:<tag>`

**为什么用 ghcr.io 而不是 Docker Hub？**
- 与 GitHub Actions 深度集成，`GITHUB_TOKEN` 自动有推送权限
- 不依赖第三方服务
- ECS 测试过能从 ghcr.io 拉取（国内网络可用）

**镜像命名规则：**
```
ghcr.io/528085390/bbs-go/auth:sha-a1b2c3d4
ghcr.io/528085390/bbs-go/auth:latest
ghcr.io/528085390/bbs-go/user:sha-a1b2c3d4
...
```

---

### 3. Kustomize（配置管理）

**目录结构：**
```
deploy/k8s/
├── base/                      # 基础配置（公共部分）
│   ├── kustomization.yaml     # 资源列表
│   ├── namespace.yaml         # 命名空间
│   ├── configmap.yaml         # 共享环境变量
│   ├── pvc.yaml               # 持久化存储
│   ├── postgres.yaml          # 数据库
│   ├── etcd.yaml              # 服务发现
│   ├── rabbitmq.yaml          # 消息队列
│   ├── user.yaml ~ gateway.yaml  # 9 个微服务
│
└── overlays/
    └── dev/                   # 开发环境覆盖层
        ├── kustomization.yaml # 镜像标签 + 补丁
        └── ingress.yaml       # 外网入口
```

**Kustomize 做了什么？**
- **Base**：定义所有资源的"模版"（通用的副本数、资源限制、探针等）
- **Overlay**：用 `patches` 和 `images` 字段覆盖 base 的值，实现环境差异化
- CI 修改 `overlays/dev/kustomization.yaml` 中的镜像标签 → ArgoCD 检测到变化 → 触发部署

**为什么用 Kustomize 而不是直接写 YAML？**
- 避免重复：9 个服务的 Deployment 结构相同，base 定义一次
- 镜像替换方便：一行 `kustomize edit set image` 就完成

---

### 4. ArgoCD（持续部署 CD）

**安装方式**：core-install + 镜像从阿里云 ACR 拉取

**核心概念：**

| 概念 | 说明 |
|---|---|
| **Application** | 一个 Git 仓库路径 → K8s 命名空间的映射 |
| **Sync** | 将 Git 中的 YAML 应用到 K8s 集群 |
| **Auto-Sync** | 自动检测 Git 变化并同步 |
| **Prune** | 删除 Git 中已移除的 K8s 资源 |
| **Self-Heal** | 有人手动改了集群，ArgoCD 会改回 Git 定义的状态 |

**同步机制：**
- 默认每 3 分钟轮询 Git 仓库
- 可通过 Webhook 实现秒级响应
- 界面可手动点击 Refresh 立即触发

**ArgoCD 工作流程：**
```
1. 轮询 GitHub 仓库 → 发现 kustomization.yaml 有更新
2. 执行 kustomize build → 生成最终的 K8s YAML
3. 对比当前集群状态与 YAML 的差异
4. 执行 kubectl apply（滚动更新）
5. 等待 Pod 就绪 → 标记同步完成
```

---

### 5. K3s / Kubernetes（容器编排）

**运行时**：Docker（k3s --docker）

**滚动更新过程（以 auth 服务为例）：**
```
秒 0:  auth-v1 Pod A (Running)
秒 1:  → auth-v2 Pod B (Pending)
秒 3:  → auth-v2 Pod B (Running, Ready)
秒 4:  → auth-v1 Pod A (Terminating)
秒 6:  → auth-v1 Pod A (Terminated)
秒 7:  → auth-v2 Pod B (Running) ✅

整个过程中，服务始终可用，没有停机时间。
```

**关键机制：**

| 机制 | 作用 |
|---|---|
| **ReplicaSet** | 保证指定数量的 Pod 始终运行 |
| **Readiness Probe** | Pod 就绪后才接入流量 |
| **Liveness Probe** | Pod 异常时自动重启 |
| **Service** | 固定的内部 DNS 名称，Pod IP 变化不影响 |
| **Ingress** | 外网统一入口，路径分发到不同 Service |
| **PVC** | 持久化存储，Pod 重启数据不丢 |

---

## 三、完整流程时间线

```
你: git push origin master
    │
    ├── 0s  GitHub 收到推送
    ├── 5s  GitHub Actions 开始运行
    │        ├── test 阶段（~30s）
    │        ├── build 阶段（~3min，9 个服务并行）
    │        └── update-config 阶段（~10s）
    │
    ├── ~4min  CI 完成，kustomization.yaml 已更新
    │
    ├── ~7min  ArgoCD 轮询到变化，开始同步（最晚 3min 内）
    │         ├── kustomize build → apply
    │         ├── K3s 开始滚动更新
    │         └── 新版本上线
    │
    └── ~8min  部署完成 ✅
```

---

## 四、回滚机制

### 自动回滚（ArgoCD）
- 在 ArgoCD UI 中点 **Sync** → 选择上一个版本
- ArgoCD 会重新应用旧版本的配置

### Git 回滚（推荐）
```bash
git revert HEAD --no-edit  # 撤销上一次提交
git push origin master      # ArgoCD 自动降级到旧版本
```

### kubectl 回滚
```bash
kubectl rollout undo deployment/auth -n temp
```

---

## 五、常见运维操作

| 操作 | 命令 |
|---|---|
| 查看 CI 状态 | https://github.com/528085390/bbs-go/actions |
| 查看 ArgoCD UI | `https://<ECS公网IP>:32446` |
| 手动触发同步 | ArgoCD UI → bbs-go → Refresh |
| 查看 Pod 日志 | `kubectl logs -n temp deployment/auth` |
| 扩缩容 | `kubectl scale deployment/auth -n temp --replicas=3` |
| 手动重启 | `kubectl rollout restart deployment/auth -n temp` |
| 查看镜像版本 | `kubectl describe deployment/auth -n temp \| grep Image` |

---

## 六、总结

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  开发者只做一件事：git push                                  │
│                                                             │
│  ┌────────┐    ┌────────┐    ┌────────┐    ┌───────────┐   │
│  │ 改代码  │ → │ git add │ → │ commit │ → │ git push  │   │
│  └────────┘    └────────┘    └────────┘    └─────┬─────┘   │
│                                                   │         │
│          ┌────────────────────────────────────────┘         │
│          ▼                                                    │
│  自动化完成：测试 → 构建 → 推送 → 更新配置 → 同步 → 上线   │
│                                                             │
│  全程无需 SSH 到服务器，无需手动 docker build                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```
