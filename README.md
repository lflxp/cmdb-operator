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

# 核心代码

```go
// Reconcile reads that state of the cluster for a CmdbService object and makes changes based on the state read
// and what is in the CmdbService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCmdbService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CmdbService")

	// Fetch the CmdbService instance
	instance := &appv1.CmdbService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		return reconcile.Result{}, err
	}

	deploy := &appsv1.Deployment{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, deploy); err != nil && errors.IsNotFound(err) {
		// 创建关联资源
		// 1. 创建Deploy
		deploy := resources.NewDeploy(instance)
		if err := r.client.Create(context.TODO(), deploy); err != nil {
			return reconcile.Result{}, err
		}
		// 2. 创建 Service
		service := resources.NewService(instance)
		if err := r.client.Create(context.TODO(), service); err != nil {
			return reconcile.Result{}, err
		}
		// 3. 关联 Annotations
		data, err := json.Marshal(instance.Spec)
		if err != nil {
			// fmt.Println("1111111111111111111", err.Error())
			return reconcile.Result{}, err
		}

		if instance.Annotations != nil {
			// fmt.Println("333333333333333333333", string(data))
			instance.Annotations["spec"] = string(data)
		} else {
			// fmt.Println("44444444444444444444444", map[string]string{"spec": string(data)})
			instance.Annotations = map[string]string{"spec": string(data)}
		}
		// 4. 更新 CmdbSerivce
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	oldspec := &appv1.CmdbService{}
	// fmt.Println("55555555555555", instance.Annotations["spec"])
	if err := json.Unmarshal([]byte(instance.Annotations["spec"]), oldspec); err != nil {
		// data, _ := json.Marshal(instance)
		// fmt.Println("222222222222222222", err.Error(), string(data))
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual(instance.Spec, oldspec) {
		// 更新关联资源
		// Deployment
		newDeploy := resources.NewDeploy(instance)
		oldDeploy := &appsv1.Deployment{}
		if err := r.client.Get(context.TODO(), request.NamespacedName, oldDeploy); err != nil {
			return reconcile.Result{}, err
		}

		oldDeploy.Spec = newDeploy.Spec
		if err := r.client.Update(context.TODO(), oldDeploy); err != nil {
			return reconcile.Result{}, err
		}

		// Service
		newService := resources.NewService(instance)
		oldService := &corev1.Service{}
		if err := r.client.Get(context.TODO(), request.NamespacedName, oldService); err != nil {
			return reconcile.Result{}, err
		}

		oldService.Spec = newService.Spec
		if err := r.client.Update(context.TODO(), oldService); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil

	// // Define a new Pod object
	// pod := newPodForCR(instance)

	// // Set CmdbService instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Check if this Pod already exists
	// found := &corev1.Pod{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	// 	err = r.client.Create(context.TODO(), pod)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Pod created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Pod already exists - don't requeue
	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	// return reconcile.Result{}, nil
}
```