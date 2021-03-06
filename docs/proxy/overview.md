# Capsule Proxy
Capsule Proxy is an add-on for the Capsule Operator.

## The problem
Kubernetes RBAC lacks the ability to list only the owned cluster-scoped resources since there are no ACL-filtered APIs. For example:

```
$ kubectl get namespaces
Error from server (Forbidden): namespaces is forbidden:
User "alice" cannot list resource "namespaces" in API group "" at the cluster scope
```

However, the user can have permissions on some namespaces

```
$ kubectl auth can-i [get|list|watch|delete] ns oil-production
yes
```

The reason, as the error message reported, is that the RBAC _list_ action is available only at Cluster-Scope and it is not granted to users without appropriate permissions.

To overcome this problem, many Kubernetes distributions introduced mirrored custom resources supported by a custom set of ACL-filtered APIs. However, this leads to radically change the user's experience of Kubernetes by introducing hard customizations that make painfull to move from one distribution to another.

With **Capsule**, we taken a different approach. As one of the key goals, we want to keep the same user's experience on all the distributions of Kubernetes. We want people to use the standard tools they already know and love and it should just work.

## How it works
This project is an add-on of the Capsule Operator, so make sure you have a working instance of Caspule before to attempt to install it. Use the `capsule-proxy` only if you want Tenant Owners to list their own Cluster-Scope resources.

The `capsule-proxy`  implements a simple reverse proxy that intercepts only specific requests to the APIs server and Capsule does all the magic behind the scenes. 

Current implementation only filter two type of requests:

* `api/v1/namespaces`
* `api/v1/nodes`
* `apis/storage.k8s.io/v1/storageclasses{/name}`
* `apis/networking.k8s.io/{v1,v1beta1}/ingressclasses{/name}`

All other requestes are proxied transparently to the APIs server, so no side-effects are expected. We're planning to add new APIs in the future, so PRs are welcome!

## Installation

The `capsule-proxy` can be deployed in standalone mode, e.g. running as a pod bridging any Kubernetes client to the APIs server.
Optionally, it can be deployed as a sidecar container in the backend of a dashboard.
Running outside a Kubernetes cluster is also viable, although a valid `KUBECONFIG` file must be provided, using the environment variable `KUBECONFIG` or the default file in `$HOME/.kube/config`.

An Helm Chart is available [here](./charts/capsule-proxy/README.md).

## Does it work with kubectl?

Yes, it works by intercepting all the requests from the `kubectl` client directed to the APIs server. It works with both users who use the TLS certificate authentication and those who use OIDC.

### Namespaces

As tenant owner `alice`, you can use `kubectl` to create some namespaces:
```
$ kubectl --context alice-oidc@mycluster create namespace oil-production
$ kubectl --context alice-oidc@mycluster create namespace oil-development
$ kubectl --context alice-oidc@mycluster create namespace gas-marketing
```

and list only those namespaces:

```
$ kubectl --context alice-oidc@mycluster get namespaces
NAME                STATUS   AGE
gas-marketing       Active   2m
oil-development     Active   2m
oil-production      Active   2m
```

### Nodes

When a Tenant defines a `.spec.nodeSelector`, the nodes matching that labels can be easily retrieved.
The annotation `capsule.clastix.io/enable-node-listing` allows the ability for the owners to retrieve the node list (useful in shared HW scenarios).

```yaml
apiVersion: capsule.clastix.io/v1alpha1
kind: Tenant
metadata:
  name: oil
  annotations:
    capsule.clastix.io/enable-node-listing: "true"
spec:
  owner:
    kind: User
    name: alice
  nodeSelector:
    kubernetes.io/hostname: capsule-gold-qwerty
```

```bash
$ kubectl --context alice-oidc@mycluster get nodes
NAME                    STATUS   ROLES    AGE   VERSION
capsule-gold-qwerty     Ready    <none>   43h   v1.19.1
```

> The following Tenant annotations allow a sort of RBAC on the operations of the nodes:
> 
> - `capsule.clastix.io/enable-node-listing`: allows listing of nodes and node retrieval
> - `capsule.clastix.io/enable-node-update`: allows the update of the node (cordoning and uncording, node tainting)
> - `capsule.clastix.io/enable-node-deletion`: allows deletion of the node

### Storage Classes

A Tenant may be limited to use a set of allowed Storage Class resources, as follows.

```yaml
apiVersion: capsule.clastix.io/v1alpha1
kind: Tenant
metadata:
  name: oil
  annotations:
    capsule.clastix.io/enable-storageclass-listing: "true"
spec:
  owner:
    kind: User
    name: alice
  storageClasses:
    allowed:
      - custom
    allowedRegex: "\\w+fs"
```

In the Kubernetes cluster we could have more Storage Class resources, some of them forbidden and non-usable by the Tenant owner.

```bash
$ kubectl --context admin@mycluster get storageclasses
NAME                 PROVISIONER              RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
cephfs               rook.io/cephfs           Delete          WaitForFirstConsumer   false                  21h
custom               custom.tls/provisioner   Delete          WaitForFirstConsumer   false                  43h
default(standard)    rancher.io/local-path    Delete          WaitForFirstConsumer   false                  43h
glusterfs            rook.io/glusterfs        Delete          WaitForFirstConsumer   false                  54m
zol                  zfs-on-linux/zfs         Delete          WaitForFirstConsumer   false                  54m
```

The expected output using `capsule-proxy` is the retrieval of the `custom` Storage Class as well the other ones matching the regex `\w+fs`.

```bash
$ kubectl --context alice-oidc@mycluster get storageclasses
NAME                 PROVISIONER              RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
cephfs               rook.io/cephfs           Delete          WaitForFirstConsumer   false                  21h
custom               custom.tls/provisioner   Delete          WaitForFirstConsumer   false                  43h
glusterfs            rook.io/glusterfs        Delete          WaitForFirstConsumer   false                  54m
```

> The following Tenant annotations allow a sort of RBAC on the Storage Class operations:
>
> - `capsule.clastix.io/enable-storageclass-listing`: allows listing of Storage Class and Storage Classes retrieval
> - `capsule.clastix.io/enable-storageclass-update`: allows the update of the Storage Class
> - `capsule.clastix.io/enable-storageclass-deletion`: allows deletion of the Storage Class

### Ingress Classes

As for Storage Class, also Ingress Class can be enforced.

```yaml
apiVersion: capsule.clastix.io/v1alpha1
kind: Tenant
metadata:
  name: oil
  annotations:
    capsule.clastix.io/enable-ingressclass-listing: "true"
spec:
  owner:
    kind: User
    name: alice
  ingressClasses:
    allowed:
      - custom
    allowedRegex: "\\w+-lb"
```

In the Kubernetes cluster we could have more Ingress Class resources, some of them forbidden and non-usable by the Tenant owner.

```bash
$ kubectl --context admin@mycluster get ingressclasses
NAME              CONTROLLER                 PARAMETERS                                      AGE
custom            example.com/custom         IngressParameters.k8s.example.com/custom        24h
external-lb       example.com/external       IngressParameters.k8s.example.com/external-lb   2s
haproxy-ingress   haproxy.tech/ingress                                                       4d
internal-lb       example.com/internal       IngressParameters.k8s.example.com/external-lb   15m
nginx             nginx.plus/ingress                                                         5d
```

The expected output using `capsule-proxy` is the retrieval of the `custom` Ingress Class as well the other ones matching the regex `\w+-lb`.

```bash
$ kubectl --context admin@mycluster get ingressclasses
NAME              CONTROLLER                 PARAMETERS                                      AGE
custom            example.com/custom         IngressParameters.k8s.example.com/custom        24h
external-lb       example.com/external       IngressParameters.k8s.example.com/external-lb   2s
internal-lb       example.com/internal       IngressParameters.k8s.example.com/internal-lb   15m
```

> The following Tenant annotations allow a sort of RBAC on the Ingress Class operations:
>
> - `capsule.clastix.io/enable-ingressclass-listing`: allows listing of Ingress Class and Ingress Classes retrieval
> - `capsule.clastix.io/enable-ingressclass-update`: allows the update of the Ingress Class
> - `capsule.clastix.io/enable-ingressclass-deletion`: allows deletion of the Ingress Class

### Storage and Ingress class required label

For Storage and Ingress Class resources, the `name` label reflecting the resource name is mandatory, otherwise filtering of resources cannot be put in place.

```yaml
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  labels:
    name: my-storage-class
  name: my-storage-class
provisioner: org.tld/my-storage-class
---
apiVersion: networking.k8s.io/v1
kind: IngressClass
metadata:
  labels:
    name: external-lb
  name: external-lb
spec:
  controller: example.com/ingress-controller
  parameters:
    apiGroup: k8s.example.com
    kind: IngressParameters
    name: external-lb
```

## Does it work with kubectl?
Yes, it works by intercepting all the requests from the `kubectl` client directed to the APIs server. It works with both users who use the TLS certificate authentication and those who use OIDC.

As tenant owner `alice`, you are able to use `kubectl` to create some namespaces:
```
$ kubectl --context alice-oidc@mycluster create namespace oil-production
$ kubectl --context alice-oidc@mycluster create namespace oil-development
$ kubectl --context alice-oidc@mycluster create namespace gas-marketing
```

and list only those namespaces:
```
$ kubectl --context alice-oidc@mycluster get namespaces
NAME                STATUS   AGE
gas-marketing       Active   2m
oil-development     Active   2m
oil-production      Active   2m
```


# What’s next
Have a fun with `capsule-proxy`:

* [Standalone Installation](./standalone.md)
* [Sidecar Installation](./sidecar.md)