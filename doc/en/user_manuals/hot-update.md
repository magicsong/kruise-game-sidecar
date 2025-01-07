## Hot Update
During the operation of the pod, it may be necessary to perform a hot update on the resources / configuration files in the pod to ensure that the pod does not restart.

### Design Architecture

![sidecar-hot-update.png](../../img/sidecar-hot-update.png)
- The hot update plugin supports triggering hot updates through semaphores.
- The configuration of the hot update plugin is shown as follows:
    - Declare the process name of the main container and the name of the semaphore to be sent in the configuration. During the hot update, the sidecar will send the specified signal to this process to trigger the main container to reread the configuration / resource files.
        - The sidecar is only responsible for sending the semaphore to the process of the main container. The main container implements reloading the configuration by itself according to the received semaphore.
- When a hot update is needed, the user accesses the port of the sidecar in this pod through an HTTP request and specifies the file path and version of the hot update in the request to trigger the hot update.
- The sidecar downloads the file from the remote end according to the hot update file path in the user's request, saves it to the specified directory of the sidecar, and this directory is mounted through emptyDir. The main container also mounts this emptyDir. In this way, the main container can obtain the new hot update file and then can be triggered to reload this file.
- After the hot update is completed, the sidecar saves the hot update result to the anno/label of the pod. And it is persisted into the specified ConfigMap (sidecar-result);
    - This ConfigMap saves the hot update results, the latest version, and the file addresses of the latest version configuration of all pods. In this way, the latest configuration can still be obtained after the pod restarts or scales.

## Usage Instructions
### GameServer 
The game server needs to have a mechanism for reloading the configuration, and the reloading can be triggered by sending a semaphore to the process of the game server container.

### Plugin Configuration
The following is the plugin configuration for hot update, which includes the following contents:
- loadPatchType: The way to trigger the hot update. Currently, the signal way is supported.
- signal:
  processName: The name of the process that triggers the hot update of the main process.
  signalName: The semaphore that needs to be sent when triggering the hot update of the main process.
- storageConfig:
    - inKube:
        - annotationKey: Which anno of the pod the hot update result is saved to.
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: sidecar-config
  namespace: kube-system
data:
  config.yaml: |
    plugins:
        - name: hot_update
          config:
            loadPatchType: signal # Hot update method
            signal:
                processName: 'nginx: master process nginx' # The process name of the main container to which the semaphore needs to be sent
                signalName: SIGHUP # Semaphore
            storageConfig: # Result saving configuration
                inKube:
                    annotationKey: sidecar.vke.volcengine.com/hot-update-result
                type: InKube # The result is saved in the k8s resource
    restartPolicy: Always
    resources:
        CPU: 100m
        Memory: 128Mi
    sidecarStartOrder: Before
```

### Example YAML of Game Server
The following provides an example of a game. In this example, by sending the SIGHUP semaphore to the 'nginx: master process nginx' process of the main container, the main container of the game server is triggered to reload the new configuration in the /var/www/html directory.
In the sidecar, an HTTP service is started and port 5000 is exposed. When a hot update is needed, access the 5000 port of the sidecar through an HTTP request and bring the version of the hot update configuration and the URL address of the configuration file in the request.
The sidecar will then pull the new configuration file according to this address and save it to the /app/downloads directory of the sidecar. And this directory will be mounted to /var/www/html of the main container, so that the main container can perceive the new configuration file.
At the same time, the sidecar will send the SIGHUP semaphore to the 'nginx: master process nginx' process of the main container to trigger the main container to reload the configuration file and complete the hot update.
- Log in to the game to view the game page after the game is deployed. At this time, it is the normal 2048 mini-game page.
- Upload the new configuration file to an address accessible by the game server. Send a curl request to the hot-update path of port 5000 of the game server to trigger the sidecar of the game server to download the configuration file and trigger the hot update.
  At this time, when accessing this game server again, you can see that the screen has been updated to V2.
```bash
curl -k -X POST -d "version=v1&url=https://xxx/v2/index.html" http://12.xxx.xx.xx:5000/hot-update 
```

```yaml
apiVersion: game.kruise.io/v1alpha1
kind: GameServerSet
metadata:
  name: hot-update-gs
  namespace: default
spec:
  replicas: 2
  gameServerTemplate:
    metadata:
      labels:
        app: hot-update-gs
    spec:
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccountName: sidecar-sa
      terminationGracePeriodSeconds: 30
      shareProcessNamespace: true
      containers:
        - image: 2048
          name: game-2048
          ports:
            - containerPort: 80
              name: game
              protocol: TCP
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          # Mount the external access address and port into /etc/ipinfo in the container
          volumeMounts:
            - mountPath: /var/www/html
              name: share-data # Persistent storage shared directory
        - image: okg-sidecar:v1
          imagePullPolicy: Always
          name: sidecar
          env:
           - name: POD_NAME
             valueFrom:
               fieldRef:
                 apiVersion: v1
                 fieldPath: metadata.name
           - name: POD_NAMESPACE
             valueFrom:
               fieldRef:
                 apiVersion: v1
                 fieldPath: metadata.namespace
          ports:
            - containerPort: 5000
              name: hot-update-port
              protocol: TCP
            - containerPort: 80
              name: game
              protocol: TCP
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /opt/sidecar
              name: sidecar-config
            - mountPath: /app/downloads
              name: share-data # Persistent storage shared directory    
      dnsPolicy: ClusterFirst
      volumes:
        - emptyDir: {}
          name: share-data
        - configMap:
            defaultMode: 420
            name: sidecar-config
          name: sidecar-config
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sidecar-result
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: sidecar-role
rules:
  - apiGroups:
      - ""
    resources:
      - pods
      - configmaps
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sidecar-sa       # Set the serviceAccount name for your pod
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: sidecar-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: sidecar-role
subjects:
  - kind: ServiceAccount
    name: sidecar-sa
    namespace: default
``` 