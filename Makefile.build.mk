export GOPROXY=https://proxy.golang.org
export CGO_ENABLED=0

SHELL := /bin/bash -o pipefail
VERSION_PACKAGE = github.com/crdant/replicated-license-enforcer/pkg/version
VERSION?=$(if $(GIT_TAG),$(GIT_TAG),alpha)
DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"`

GIT_TREE = $(shell git rev-parse --is-inside-work-tree 2>/dev/null)
ifneq "$(GIT_TREE)" ""
define GIT_UPDATE_INDEX_CMD
git update-index --assume-unchanged
endef
define GIT_SHA
`git rev-parse HEAD`
endef
else
define GIT_UPDATE_INDEX_CMD
echo "Not a git repo, skipping git update-index"
endef
define GIT_SHA
""
endef
endif

ifeq ("$(DEBUG_REPLICATED)", "1")
define LDFLAGS
-ldflags "\
	-X ${VERSION_PACKAGE}.Version=${VERSION} \
	-X ${VERSION_PACKAGE}.GitSHA=${GIT_SHA} \
	-X ${VERSION_PACKAGE}.BuildTime=${DATE} \
"
endef
define GCFLAGS
-gcflags="all=-N -l"
endef
else
define LDFLAGS
-ldflags "\
	-s -w \
	-X ${VERSION_PACKAGE}.Version=${VERSION} \
	-X ${VERSION_PACKAGE}.GitSHA=${GIT_SHA} \
	-X ${VERSION_PACKAGE}.BuildTime=${DATE} \
"
endef
endif

BUILDFLAGS = -tags='netgo containers_image_ostree_stub exclude_graphdriver_devicemapper exclude_graphdriver_btrfs containers_image_openpgp' -installsuffix netgo
TEST_BUILDFLAGS = -tags='testing netgo containers_image_ostree_stub exclude_graphdriver_devicemapper exclude_graphdriver_btrfs containers_image_openpgp' -installsuffix netgo
