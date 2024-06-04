# Replicated License Enforcer

Implements a simple command that validates if a Replicated license is expired.
Use it in an container to prevent your application from starting after a
license has expired.

[![Enforcing the Replicated License at Runtime](https://cdn.loom.com/sessions/thumbnails/81f608f80ca1493dbed01584d82fb5b9-with-play.gif)](https://www.loom.com/share/81f608f80ca1493dbed01584d82fb5b9)

## Background

Runtime license enforcement is a common theme when I discuss the Replicated
Platform with software vendors. We deliberately focus our enforcement efforts
on distribution and installation. I think that makes sense, since we're not
sure how you might want to respond when the license isn't valid. Instead, we
give you the tools you need to write your own enforcement code. The
[Replicated SDK](https://docs.replicated.com/reference/replicated-sdk-apis)
gives you access to the license and all its fields. You can easily write cod3
to access the expiration date and any of your own entitlements. 

But something kept gnawing at me. Was there a way I could create something
more general that multiple vendors could use. Not something for Replicated to
ship (necessarily), but something that folks could use directly or as
inspiration. Turns out, there was something I could do.


## Usage

Regardless of how you're using the code, there's a basic assumption that you
are [using the Replicated
SDK](https://docs.replicated.com/vendor/replicated-sdk-overview). If you're
not, you need to create a dependency on it for your Helm chart.

```
dependencies:
- name: replicated
  repository: oci://registry.replicated.com/library
  version: 1.0.0-beta.20
```

The version might have changed since I last updated this README, so take a
look at the [latest releases](https://github.com/replicatedhq/replicated-sdk/releases).

### As an init container

The most direct way to use this code is to add an init container to
one or more workloads in your application that uses my image directly. The
init container will stop your pod from running until the Replicated license is
valid. 

Add the following code to your manifest:

```
initContainers:
- name: license-check
  image: ghcr.io/crdant/replicated-license-enforcer:latest
  env:
    - name: REPLICATED_SDK_ENDPOINT
      value: http://replicated:3000
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
    - name: POD_UID
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.uid
  imagePullPolicy: IfNotPresent
```

A couple of things to be aware of:

1. You must be using the Replicated SDK and have an active Replicated
   subscription or trial.
2. I recommend pulling the image through the Replicated proxy using your the
   custom domain you've set up for it. You'll need the right image pull secret
   if you do that.
3. The code creates Kubernetes events, so you need appropriate RBAC. You'll
   want to create your a service account and assign it an appropriate role.
   See [`examples/rbac.yaml`](./examples/rbac.yaml).

### As a sidecar

The init container approach is nice because it stops the pod from running
without a valid license. It's constrained, though. Once the pod is running the
license isn't checked again until it restarts. What happens if the license
expires between restarts?

Running the container as a sidecar addresses this scenario. The sidecar will
not stop other containers in the pod, but it will record Kubernetes events to
show the license is no longer valid. Valid license checks will not emit new
events unless the expiration date has changed, and the code will not check
again until a set duration has passed. You set the duration with the
`--recheck` flag passed to the `enforcer` command.

Use the following to run as a sidecar:

```
initContainers:
- name: license-check
  image: ghcr.io/crdant/replicated-license-enforcer:latest
  command: [ "/enforcer" ]
  args: ["--recheck", "4h" ]    # any valid Go duration
  env:
    - name: REPLICATED_SDK_ENDPOINT
      value: http://replicated:3000
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
    - name: POD_UID
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.uid
  imagePullPolicy: IfNotPresent
```

### In your own code

The core packages in this repository are re-usable in your own license
enforcement code. 

#### `enforce` Package

The main enforcement code is the package `enforce`. The `enforce` command is a
simple CLI wrapper around that package, so you can 

To validate the license as your code starts up (just like the init container):

```go
package main

import (
    "github.com/crdant/replicated-license-enforcer/pkg/enforce"
)

func main() {
    // Check license before we start
    enforcer := enforce.DefaultEnforcer()
	err := enforcer.Validate()
	if err != nil {
		log.Error("Error checking license validity", "error", err)
		os.Exit(1)
	}

   // rest of your code here 
}
```

Your code won't run until the valid license is in place, and you'll see the
license valid/expired events associated with the pod in Kubernetes.

If you want to continue monitor the license as your application runs, with the
same valid/invalid events generated, use `Monitor` after you initial
validation completes validation.

```go
package main

import (
    "github.com/crdant/replicated-license-enforcer/pkg/enforce"
)

func main() {
    // Check license before we start
    enforcer := enforce.DefaultEnforcer()
	err := enforcer.Validate()
	if err != nil {
		log.Error("Error checking license validity", "error", err)
		os.Exit(1)
	}

    // set the interval for rechecking the license, license expire on a given
    // day at 00:00:00 UTC, but checking more often can help with picking up 
    // a license that has been extended
    recheckInterval := 4h
    enforcer.Monitor(recheckInterval)

    // rest of your code here 
}
```

#### `client` Package

The `client` package is a client for the Replicated SDK that focuses on the
fields that are most useful for license enforcement. The client is simplified
for this purpose, but may evolve into a more complete client over time.

#### `event` Package

A purpose-built package for recording Kubernetes events related to the state
of the Replicated license.

## Limitations

1. This code is provided _AS IS_ and is not supported by Replicated.

2. Your customer can circumvent this by editing Kubernetes manifests. I don't
   consider that too much of a limitation because you still have legal
   remedies. It'll even help your case that the worked to circumvent your
   enforcement.
