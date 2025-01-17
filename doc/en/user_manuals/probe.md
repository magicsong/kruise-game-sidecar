## Probe
OpenKruiseGame provides the ability to customize Probes through PodProbeMarker and returns the results to the GameServer. However, this ability cannot be used in the serverless scenario.
This sidecar designs a PodProbe plugin to implement Probes and solve the problems in the serverless scenario.

### Design Architecture
![pod-probe](../../img/sidecar-pod-probe.png)
- Set the configurations required by the podProbe plugin through the configured ConfigMap. The configurations of the podProbe plugin mainly include the following parts:
  - url: Detection address. Currently, HTTP interface detection is supported.
  - method: HTTP request method.
  - timeout: HTTP request timeout period.
  - expectedStatusCode: Expected HTTP status code.
  - storageConfig: Detection result storage configuration, which can store the results in the opsState of the GameServer.
- Conduct cyclic detection according to the set detection address, and modify the status of spec.opsState of the GameServer according to the detection results. The specific rules are defined in the configured ConfigMap.

## Usage Instructions
### GameServer 
Since this plugin uses HTTP to detect and obtain information about the game server, the information that needs to be detected needs to be exposed to a certain port through an HTTP service so that the sidecar can conduct the detection.

### Plugin Configuration
The configuration information required by the plugin needs to be written into the ConfigMap for the sidecar to use. This ConfigMap is mounted to the sidecar of the GameServerSet in a mounted manner.
The following is an example of the ConfigMap for service quality detection, which includes information such as the HTTP detection method, status code, detection address, detection start delay, detection interval, and detection timeout period.
In addition, the detection result saving configuration is also defined. The markerPolices declare the detection rules. For example, when the detection result is WaitToBeDeleted, the spec.opsState in the GameServer will be set to WaitToBeDeleted.
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: sidecar-config
  namespace: sidecar
data:
  config.yaml: |
    plugins:
      - name: http_probe
        config:
          endpoints:
            - expectedStatusCode: 200 ## The status code representing successful detection
              method: GET # HTTP request method
              storageConfig: # Result saving configuration
                inKube:
                  markerPolices: # Result saving rules: When the request result is WaitToBeDeleted, set spec.opsState in the GameServer to WaitToBeDeleted
                    - gameServerOpsState: WaitToBeDeleted
                      state: WaitToBeDeleted 
                    - gameServerOpsState: Allocated
                      state: Allocated
                    - gameServerOpsState: None
                      state: None
                type: InKube
              timeout: 30  # Detection timeout period in seconds
              url: http://localhost:8080 # HTTP detection address
          startDelaySeconds: 10 # Detection start delay in seconds
          probeIntervalSeconds: 5 # Detection interval in seconds
        bootorder: 1
    restartPolicy: Always
    resources:
      CPU: 100m
      Memory: 128Mi
    sidecarStartOrder: Before ## The startup order of the Sidecar, whether it is after or before the main container
```

### Game Server YAML Example
The sidecar can be used by configuring it into the YAML of the GameServerSet. The following provides an example of a GameServerSet, which includes the sidecar container settings and mounts the ConfigMap to the sidecar container.
After the GameServerSet is started, the sidecar will conduct cyclic detection on the status of the game server according to the configuration and set the results to the spec.opsState of the GameServer.
In this example, a TankWar game is provided, and whether there are players is exposed through the HTTP service on port 8080. The sidecar can detect whether there are players on the game server through port 8080.
- In the initial state, when the detection result is None, the sidecar will set the spec.opsState of the GameServer to None.
- When accessing the address to log in to the game server, the sidecar will set the spec.opsState of the GameServer to Allocated.
- When players on the game server exit, the sidecar will set the spec.opsState of the GameServer to WaitToBeDeleted.
```yaml
apiVersion: game.kruise.io/v1alpha1
kind: GameServerSet
metadata:
  name: hot-update-gs
  namespace: sidecar
spec:
  replicas: 2
  updateStrategy:
    rollingUpdate:
      podUpdatePolicy: InPlaceIfPossible
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
        - image: tank:latest
          name: tankwar
          env:
            - name: DEBUG
              value: ts-mp:*
          resources:
            limits:
              cpu: 1500m
              memory: 3000Mi
            requests:
              cpu: 1500m
              memory: 3000Mi
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
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /opt/sidecar
              name: sidecar-config
      dnsPolicy: ClusterFirst
      volumes:        
        - configMap:
            defaultMode: 420
            name: sidecar-config
            namespace: sidecar
          name: sidecar-config
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
  - apiGroups:
      - "game.kruise.io"
    resources:
      - '*'
    verbs:
      - patch
      - update
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sidecar-sa       # Set the serviceAccount name for your pod
  namespace: sidecar
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
    namespace: sidecar
``` 