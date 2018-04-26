# authserver
This [helm](https://github.com/kubernetes/helm) chart provides
[kubernetes](http://kubernetes.io) manifests for running
[authserver](https://github.com/dictyBase/authserver).

## Configuration

The following tables lists the configurable parameters of the **chado-sqitch** chart and their default values.

| Parameter                 | Description                           | Default                                                  |
| --------------------------|---------------------------------------| --------------------------------------------|
| `image.repository`        | authserver image                      | `dictybase/authserver`                      |
| `image.tag`               | image tag                             | `3.0.0`                                     |
| `image.pullPolicy`        | Image pull policy                     | `IfNotPresent`                              |
| `service.name`            | Name of the service.                  | `authserver`                                |
| `service.type`            | Type of service.                      | `NodePort`                                  |
| `service.port`            | Port of service.                      | `9999`                                      |
| `resources`               | CPU/Memory resource requests/limits   |  `nil`                                      |
| `nodeSelector`            | Node labels for pod assignment        |  `nil`                                      |
| `publicKey`               | Public key(string) read from file     |  `Have to be set from command line`         |                        |
| `privateKey`              | Public key(string) read from file     |  `Have to be set from command line`         |                                          |
| `configFile`              | Client secrets read from file         |  `Have to be set from command line`         |                                          |

### Required configurations
The three parameters `publicKey`, `privateKey` and `configFile` have to be set
from the command line. The keyfiles(`publicKey` and `privateKey`) can be
generated following the instructions
[here](https://github.com/dictyBase/authserver#generate-keys). The
configuration file(`configFile`) can be created following the direction
[here](https://github.com/dictyBase/authserver#create-configuration-file)

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. 

```
helm install --name my-release --set publicKey=$(base64 -w0 publicKeyfile) 
    --set privateKey=$(base64 -w0 privateKeyfile) --set configFile=$(base64 -w0 configFile) authserver
```

>Helm for security reason does not allow to include arbitrary file outside of
>charts folder. So, to read any file outside of the package, particularly
>released package, the content of the file should be read as base64 encoded
>string in a single line(-w0) and passed to the(--set
>option) chart.


