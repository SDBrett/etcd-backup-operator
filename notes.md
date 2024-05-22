Notes and planning<-- omit in toc -->

- [Bare minimum for testing](#bare-minimum-for-testing)
  - [Functionality](#functionality)
  - [Resources](#resources)
  - [Inputs](#inputs)
  - [Process](#process)
- [Making useable](#making-useable)
  - [CR Configuration Items](#cr-configuration-items)


## Bare minimum for testing

What is required to test raw functionality of backing up ETCD

### Functionality
- Backup job runs in same namespace as controller
- All resources are in same namespace as controller
- Backup to PVC

### Resources
- CronJob
- ConfigMap with backup script (easier to update CM as CronJob update can be flakey)
- PVC to store backup content
- Service account and RBAC for backup job

### Inputs
No input required from CR to test raw functionality.

### Process
1. Reconcile EtcdBackup CR
2. Create ConfigMap containing backup script
3. Create RBAC resources
4. Create CronJob

## Making useable

Improvements to make operator actually usable

### CR Configuration Items
- [ ] Namespace for ETCD backup resources
- [ ] Optional storage backend types
  - [ ] PVC
  - [ ] S3
- [ ] Encrypt Backup
- [ ] Node selection
- [ ] Tolerations
- [ ] User provided RBAC resources
- [ ] Status
- [ ] Event reporting
- [ ] ETCD settings
  - [ ] Certificate paths
    - [ ] cert
    - [ ] ca cert
    - [ ] key
  - [ ] endpoint IP (default to 127.0.0.1)
  - [ ] port (default to 2379)
  - [ ] Configure container image

