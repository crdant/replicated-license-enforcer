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

### As an init container

You can use this code. The most direct way is to add an init container to
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

### In your own code

There are a couple of useful packages you can take advantage of if you'd
rather incorporate enforcement into your own code. You can use the `client`
package as a lightweight client for the needed parts of the Replicated SDK
## Limitations

1. This code is provided _AS IS_ and is not supported by Replicated.

2. Your customer can circumvent this by editing Kubernetes manifests. I don't
   consider that too much of a limitation because you still have legal
   remedies. It'll even help your case that the worked to circumvent your
   enforcement.
