# 命令合集

* operator-sdk 
* operator-sdk version 查看版本
* operator-sdk new cmdbdemo 新建项目
* operator-sdk add api --api-version=app.example.com/v1 --kind=CmdbService 定义新的api CRD
* operator-sdk add controller --api-version=app.example.com/v1 --kind=CmdbService 定义CRD的执行逻辑
* operator-sdk generate k8s 自动生成代码

# 部署

* operator-sdk build cnych/opdemo # 应用打包成docker镜像
* docker push cnych/opdemo # 镜像构建成功后，推送到 docker hub
* sed -i 's|REPLACE_IMAGE|cnych/opdemo|g' deploy/operator.yaml # 修改镜像地址更新Operator的资源清单
* kubectl create -f deploy/service_account.yaml # 创建对应的 RBAC 的对象 Setup Service Account
* kubectl create -f deploy/role.yaml # Setup RBAC
* kubectl create -f deploy/role_binding.yaml S# etup RBAC
* kubectl apply -f deploy/crds/app_v1_appservice_crd.yaml # Setup the CRD
* kubectl create -f deploy/operator.yaml # Deploy the Operator
* kubectl create -f deploy/crds/app_v1_appservice_cr.yaml # 业务逻辑部署

# 清理

```sh
$ kubectl delete -f deploy/crds/app_v1_appservice_cr.yaml
$ kubectl delete -f deploy/operator.yaml
$ kubectl delete -f deploy/role.yaml
$ kubectl delete -f deploy/role_binding.yaml
$ kubectl delete -f deploy/service_account.yaml
$ kubectl delete -f deploy/crds/app_v1_appservice_crd.yaml
```