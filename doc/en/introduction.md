## kruise-game-sidecar
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/openkruise/kruise-game)](https://goreportcard.com/report/github.com/openkruise/kruise)
[![codecov](https://codecov.io/gh/openkruise/kruise-game/branch/master/graph/badge.svg)](https://codecov.io/gh/openkruise/kruise-game)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](./CODE_OF_CONDUCT.md)

English | [Chinese](../中文/说明文档.md)
### Overview
Currently, Kruise provides the capabilities of service quality detection and hot update. However, these capabilities are implemented in Kruise-daemon, and in the serverless scenario, Kruise-daemon cannot be used normally.
Therefore, it is considered to implement the enhanced functions in the serverless game scenario through the sidecar mechanism. This sidecar contains multiple plugins that can solve different problems in the game scenario. It can be used by integrating this sidecar into the pod.
In addition, the sidecar provides a complete and extensible plugin management mechanism, allowing you to easily add the plugins you need according to your own requirements.

### Sidecar Design Architecture
![sidecar-struct](../img/sidecar-struct.png)
1. Create a ConfigMap to store the configurations required by each plugin in the sidecar. Each plugin can customize its own required configuration and store it in this ConfigMap.
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
   name: kidecar-config
   namespace: sidecar
data:
   config.yaml: |
      plugins:
        - name: hot_update
          config:
            loadPatchType: signal # Trigger hot update through semaphore
            signal:
              processName: 'nginx: master process nginx' ## The name of the process that triggers the hot update of the main process
              signalName: SIGHUP ## The semaphore that needs to be sent when triggering the hot update of the main process
            storageConfig:
              inKube:
                annotationKey: sidecar.vke.volcengine.com/hot-update-result # Save the hot update result to which anno of the pod
              type: InKube
          bootorder: 0
      restartpolicy: Always
      resources:
        CPU: 100m
        Memory: 128Mi
      sidecarstartorder: Before ## The startup order of the Sidecar, whether it is before or after the main container

```
3. The sidecar will start and run each plugin in turn according to the plugins configured in the ConfigMap.
4. After the plugin runs, the results can be persisted. There are some preset result saving mechanisms in the sidecar, which can save the results to the anno/label of the pod or the specified location of the custom CRD.
   It can be set in the storageConfig of the plugin configuration in the ConfigMap. In addition, plugin developers can also implement result persistence by themselves and save the plugin results to the expected location.

### Supported Plugins
Currently, two plugins are preset, namely the hot update plugin and the service quality detection plugin.
* Hot update plugin: Supports the hot update of pods and supports triggering hot update through semaphore.
* Service quality detection plugin: Supports the service quality detection of game servers and supports detecting the service quality of pods through HTTP.

### Next Steps
* View [probe](./user_manuals/probe.md) to use the service quality probe plugin.
* View [Hot Update](./user_manuals/probe.md) to use the hot update plugin. 

