# Meet the multi-tenancy benchmark MTB
Actually, there's no yet a real standard for the multi-tenancy model in Kubernetes, although the [SIG multi-tenancy group](https://github.com/kubernetes-sigs/multi-tenancy) is working on that. SIG multi-tenancy drafted a generic validation schema appliable to generic multi-tenancy projects. Multi-Tenancy Benchmarks [MTB](https://github.com/kubernetes-sigs/multi-tenancy/tree/master/benchmarks) are guidelines for multi-tenant configuration of Kubernetes clusters. Capsule is an open source multi-tenacy operator and we decided to meet the requirements of MTB.

> N.B. At time of writing, the MTB are in development and not ready for usage. Strictly speaking, we do not claim an official conformance to MTB, but just to adhere to the multi-tenancy requirements and best practices promoted by MTB.

|Benchmark     |MTB ID  |MTB Profile|Capsule Version|Conformance|Notes  |
|--------------|--------|-----------|---------------|-----------|-------|
|Block access to cluster resources|MTB-PL1-CC-CPI-1|L1|v0.1.0|✓|---|
|Block access to Multitenant Resources|MTB-PL1-BC-CPI-2|L1|v0.1.0|✓|---|
|Block access to other tenant resources|n/a|L1|v0.1.0|✓|MTB draft|
|Block add capabilities|MTB-PL1-BC-CPI-3|L1|v0.1.0|✓|---|
|Require always imagePullPolicy|MTB-PL1-CC-DI-1|L1|v0.1.0|✓|---|
|Require run as non-root user|MTB-PL1-BC-CPI-4|L1|v0.1.0|✓|---|
|Block privileged containers|MTB-PL1-BC-CPI-5|L1|v0.1.0|✓|---|
|Block privilege escalation|MTB-PL1-BC-CPI-6|L1|v0.1.0|✓|---|
|Configure namespace resource quotas|MTB-PL1-CC-FNS-1|L1|v0.1.0|✓|---|
|Configure namespace object limits|MTB-PL1-CC-FNS-2|L1|v0.1.0|✓|---|
|Block use of host path volumes|MTB-PL1-BC-HI-1|L1|v0.1.0|✓|---|
|Block use of NodePort services|MTB-PL1-BC-HI-1|L1|v0.1.0|✓|---|
|Block use of host networking and ports|MTB-PL1-BC-HI-3|L1|v0.1.0|✓|---|
|Block use of host PID|MTB-PL1-BC-HI-4|L1|v0.1.0|✓|---|
|Block use of host IPC|MTB-PL1-BC-HI-5|L1|v0.1.0|✓|---|
|Block modification of resource quotas|MTB-PL1-CC-TI-1|L1|v0.1.0|✓|---|
|Require PersistentVolumeClaim for storage|n/a|L1|v0.1.0|✓|MTB draft|
|Require PV reclaim policy of delete|n/a|L1|v0.1.0|✓|MTB draft|
|Block use of existing PVs|n/a|L1|v0.1.0|✓|MTB draft|
|Block network access across tenant namespaces|n/a|L1|v0.1.0|✓|MTB draft|
|Allow self-service management of Network Policies|MTB-PL2-BC-OPS-4|L2|v0.1.0|✓|---|
|Allow self-service management of Roles|n/a|L2|v0.1.0|✓|MTB draft|
|Allow self-service management of Roles Bindings|n/a|L2|v0.1.0|✓|MTB draft|

