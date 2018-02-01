Infinibox Terraform Provider
==================

Disclaimer
----------
This project is neither supported nor maintained by Infinidat.

- Infinidat Website: https://www.infinidat.com
- Terraform Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Contents
------
1. [Requirements](#requirements) - lists the requirements for building the provider
2. [Building The Provider](#building-the-provider) - lists the steps for building the provider
3. [Using The Provider](#using-the-provider) - details how to use the provider
4. [Developing The Provider](#developing-the-provider) - steps for contributing back to the provider

Requirements
------------

-    [Terraform](https://www.terraform.io/downloads.html) 0.11.x
-    [Go](https://golang.org/doc/install) 1.9+ (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/devnal/terraform-provider-ibox`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:devnal/terraform-provider-ibox
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/devnal/terraform-provider-ibox
$ make build
```

Using The Provider
----------------------

### Infinibox Resources List and examples

1. [Provider](#provider)
2. [Pool](#pool)
3. [Volume](#volume)
4. [Host](#host)
5. [HostCluster](#hostcluster)
6. [Lun](#lun)

### Provider


Username, Password and Hostname are required for provider configuration.

_Example_
```hcl
provider "ibox" {
  hostname = "ibox630"
  username = "admin"
  password = "123456"
}
```

### Pool

[Pool Api Docs](https://ibox630/apidoc/#PoolResource)

Pool resource has to be configured with minimal physical capacity of 1TB, virtual capacity allows over provisioning.
Capacity can be increased or decreased. SSD read cache and compression can be enabled/disabled for this resource.

_Example_
```hcl
resource "ibox_pool" "my-pool" {
  name = "my-pool-test"
  physical_capacity = "1100000000000"
  virtual_capacity = "3000000000000"
  physical_capacity_critical = 95
  physical_capacity_warning = 89
  ssd_enabled = true
  compression_enabled = true
}
```

### Volume

[Volume Api Docs](https://ibox630/apidoc/#VolumeResource)

Volume resource has to be configured with minimal size of 1GB.
Capacity can only be increased or decreased.
Volume can be provisioned as THIN or THICK. 
Volume must be created in one of the pools.

_Example_
```hcl
resource "ibox_volume" "my-volume" {
  name = "my-volume-test"
  pool_id = "${ibox_pool.my-pool.id}"
  size = 20000000000
  provtype = "THIN"
  ssd_enabled = true
  compression_enabled = true
}
```

### Host

[Host Api Docs](https://ibox630/apidoc/#HostResource)

This type of resource creates Host object, for ISCSI host authentication can be added.
List of FC or/and ISCSI ports can be added during creation or updated later.

_Example_
```hcl
resource "ibox_host" "my-host" {
  name = "my-host-test"
}
```

_Example Host with ISCSI CHAP AUTH_
```hcl
resource "ibox_host" "my-host" {
  name = "my-host-test"
  security_method = "MUTUAL_CHAP"
  security_chap_inbound_username = "hostuser987"
  security_chap_inbound_secret = "hostsecret987654321"
  security_chap_outbound_username = "hostuser98732"
  security_chap_outbound_secret = "hostsecret987654321dasdasdsa"
}
```

_Example Host with ISCSI MUTUAL CHAP AUTH and FC Ports_
```hcl
resource "ibox_host" "my-host7" {
  name = "my-host-test7"
  security_method = "MUTUAL_CHAP"
  security_chap_inbound_username = "h"
  security_chap_inbound_secret = "hostsecret987654321dasdasdsaddd"
  security_chap_outbound_username = "hostuser98732"
  security_chap_outbound_secret = "hostsecret987654321dasdasdsa"
  ports = [

      {
        type = "FC"
        address = "21100024ff913bff"
      },
      {
        type = "FC"
        address = "22100024ff913bff"
      },      
  ]
}
```

### Host_Cluster

[Host Cluster Api Docs](https://ibox630/apidoc/#HostClusterResource)

Host cluster resource is a collection of host resources that are sharing same mapped volumes.
List of hosts can be created during creation or updated later.
Changing host list is a dangerous operation and must be performed carefully.

_Example_
```hcl
resource "ibox_host_cluster" "my-host-cluster" {
  name = "my-host-cluster"
  hosts = [
    "${ibox_host.my-host.id}",
    "${ibox_host.my-host2.id}",
    "${ibox_host.my-host3.id}",
  ]
}
```

### Lun

[Adding LUN to Host Cluster Api Docs](https://ibox630/apidoc/#addALunToACluster)

LUN resource can be added or removed from Host or Host Cluster resources.
If needed specific LUN ID can be defined e.g. 20

_Example_
```hcl
resource "ibox_lun" "my-host-lun-3" {
  volume_id = "${ibox_volume.my-volume3.id}"
  host_id = "${ibox_host.my-host6.id}"
  lun = 20
}
```
