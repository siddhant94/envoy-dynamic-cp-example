# envoy-dynamic-cp-example
An example containing bootstrap envoy config which has listeners and clusters updated dynamically via file. K8s Deployments for envoy as a sidecar proxy with tcp_proxy, http_connection_manager etc filters.  

Minikube (K8s) Deployments  

There are 2 config maps. The proxy-config contains the bootstrap envoy config, whereas the proxy-dynamic-config contains the lds and cds configs.
```
kubectl get cm
NAME                   DATA   AGE
proxy-config           1      64m
proxy-dynamic-config   2      63m
```

Commands to create these,  
`k create cm proxy-config --from-file envoy.yaml`  
`k create cm proxy-dynamic-config --from-file=cds.yaml=./envoy_clusters/cds_mysql.yaml --from-file=lds.yaml=./envoy_listeners/lds_mysql.yaml`  


For the deployments,  
`k apply -f ./deployment/mysql-svc-deployment.yaml` & `k apply -f ./deployment/http-svc-deployment.yaml` would deploy the mysql and http svc along with envoy proxy as sidecar. Iptable rules ensure that all traffic to the service port of http and mysql are redirected to envoy proxy which then forwards the request to respective service.  

These services are exposed as LoadBalancer type and hence we need to run `minikube tunnel` to obtain an "external" ip.  
For mysql service, we can leverage a mysql client to generate traffic.  
To run a mysql client which will connect via 3306 port to our created mysql service, we need to run  
`❯ kubectl run -it --rm --image=mysql:5.5 --restart=Never mysql-client -- mysql -h 1.2.3.4 -P 3306 -u root`. This would give us an client shell to mysql where we can run queries like `CREATE database test;`. We are connecting via port 3306 but the iptables ensure that very call to this port goes through our envoy. We can check the filter stats to ensure this.  
Check out envoy proxy's sandbox example https://www.envoyproxy.io/docs/envoy/v1.17.0/start/sandboxes/mysql from which this example is inspired.
