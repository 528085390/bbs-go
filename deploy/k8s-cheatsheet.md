# K8s / K3s 常用命令速查表

> 项目：bbs-go | 命名空间：`temp`（以下所有命令都在 `-n temp` 下操作）

---

## 一、Pod 相关

```bash
# 列出所有 Pod
kubectl get pods -n temp

# 带更多信息（IP、节点）
kubectl get pods -n temp -o wide

# 持续 watch 变化
kubectl get pods -n temp -w

# 查看 Pod 详细信息（事件、状态）
kubectl describe pod <pod-name> -n temp

# 查看 Pod 日志（实时 tail）
kubectl logs -f <pod-name> -n temp

# 查看 Deployment 级别日志（自动选当前 Pod）
kubectl logs -f deployment/auth -n temp

# 只看最后 50 行
kubectl logs --tail=50 deployment/auth -n temp

# 进入 Pod 内部（调试用）
kubectl exec -it <pod-name> -n temp -- sh

# 删除 Pod（Deployment 会自动重建）
kubectl delete pod <pod-name> -n temp

# 删除所有异常 Pod
kubectl delete pod -n temp --field-selector=status.phase=Failed
```

---

## 二、Deployment 相关

```bash
# 列出 Deployment
kubectl get deployments -n temp

# 查看详情
kubectl describe deployment auth -n temp

# 扩容/缩容
kubectl scale deployment/auth -n temp --replicas=3
kubectl scale deployment/auth -n temp --replicas=1

# 滚动更新（改镜像版本）
kubectl set image deployment/auth -n temp auth=temp-auth:v2

# 查看滚动更新状态
kubectl rollout status deployment/auth -n temp

# 查看历史版本
kubectl rollout history deployment/auth -n temp

# 回滚到上一版本
kubectl rollout undo deployment/auth -n temp

# 回滚到指定版本
kubectl rollout undo deployment/auth -n temp --to-revision=2

# 重启（重新拉取镜像、分配 IP）
kubectl rollout restart deployment/auth -n temp

# 暂停/恢复滚动更新
kubectl rollout pause deployment/auth -n temp
kubectl rollout resume deployment/auth -n temp
```

---

## 三、Service 相关

```bash
# 列出所有 Service
kubectl get services -n temp

# 查看 Service 详情（ClusterIP、Endpoints）
kubectl describe service auth-svc -n temp

# 临时将 Service 暴露到本机端口（调试用）
kubectl port-forward service/auth-svc -n temp 8001:8001
# 之后就可以 curl http://localhost:8001 访问了
```

---

## 四、Ingress 相关

```bash
# 列出 Ingress
kubectl get ingress -n temp

# 查看详情
kubectl describe ingress gateway-ingress -n temp

# 查看 Traefik Ingress Controller 日志
kubectl logs -n kube-system -l app.kubernetes.io/name=traefik --tail=50
```

---

## 五、配置（ConfigMap / Secret）

```bash
# 列出 ConfigMap
kubectl get configmaps -n temp

# 查看内容
kubectl get configmap bbs-config -n temp -o yaml

# 更新 ConfigMap（改完需要重启 Pod 生效）
kubectl edit configmap bbs-config -n temp

# ConfigMap 改完重启所有引用它的 Deployment
kubectl rollout restart deployment -n temp -l project=bbs-go
```

---

## 六、存储（PVC / PV）

```bash
# 列出 PVC
kubectl get pvc -n temp

# 查看 PVC 详情
kubectl describe pvc postgres-pvc -n temp

# 查看 PV（集群级，不需要 -n）
kubectl get pv

# 检查 PVC 挂载到 Pod 的情况
kubectl describe pod postgres-xxxxx -n temp | grep -A5 Volumes
```

---

## 七、集群管理

```bash
# 查看节点
kubectl get nodes

# 查看集群信息
kubectl cluster-info

# 查看所有命名空间的 Pod
kubectl get pods --all-namespaces

# 查看节点资源使用（需要 metrics-server）
kubectl top nodes
kubectl top pods -n temp
```

---

## 八、k3s 特有

```bash
# 查看 k3s 服务状态
systemctl status k3s

# 重启 k3s
systemctl restart k3s

# 卸载 k3s
/usr/local/bin/k3s-uninstall.sh

# k3s containerd 导入镜像
docker save temp-auth:latest | k3s ctr images import -

# k3s containerd 列出已导入镜像
k3s ctr images list | grep temp
```

---

## 九、故障排查

```bash
# 查看 Pod 最近事件（为什么 Pending / CrashLoopBackOff）
kubectl describe pod <pod-name> -n temp

# 查看容器日志
kubectl logs <pod-name> -n temp --tail=100

# 查看前一个容器实例的日志（CrashLoopBackOff 时有用）
kubectl logs <pod-name> -n temp --previous

# 检查所有资源
kubectl get all -n temp

# 查看 etcd 健康状况（Pod 内执行）
kubectl exec -n temp deploy/etcd -- etcdctl endpoint health

# 查看 PostgreSQL 是否正常
kubectl exec -n temp deploy/postgres -- pg_isready -U postgres
```

---

## 十、清理资源

```bash
# 删除整个命名空间及其所有资源（慎重！）
kubectl delete namespace temp

# 删除单个 Deployment
kubectl delete deployment auth -n temp

# 删除所有异常 Pod
kubectl delete pod -n temp --field-selector=status.phase=Failed

# 删除 PVC（会同时删除数据）
kubectl delete pvc postgres-pvc -n temp
```
