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
gives you access to the license and all its fields. You can easily write field
to access the expiration date and any of your own entitlements. 

But something kept gnawing at me. Was there a way I could create something
more general that multiple vendors could use. Not something for Replicated to
ship (necessarily), but something that folks could use directly or as
inspiration. Turns out, there was something I could do.

## Limitations

1. This code is provided _AS IS_ and is not supported by Replicated.

2. Your customer can circumvent this by editing Kubernetes manifests. I don't
   consider that too much of a limitation because you still have legal
   remedies. It'll even help your case that the worked to circumvent your
   enforcement.
