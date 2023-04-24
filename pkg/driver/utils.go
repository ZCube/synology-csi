/*
Copyright 2021 Synology Inc.

Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/mount-utils"
	"k8s.io/utils/exec"
)

func ParseEndpoint(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") || strings.HasPrefix(strings.ToLower(ep), "tcp://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
	}
	return "", "", fmt.Errorf("Invalid endpoint: %v", ep)
}

func NewControllerServer(d *Driver) *controllerServer {
	return &controllerServer{
		Driver:     d,
		dsmService: d.DsmService,
	}
}

func NewNodeServer(d *Driver) *nodeServer {
	return &nodeServer{
		Driver:     d,
		dsmService: d.DsmService,
		Mounter: &mount.SafeFormatAndMount{
			Interface: mount.New(""),
			Exec:      exec.New(),
		},
		Initiator: &initiatorDriver{
			chapUser:     "",
			chapPassword: "",
		},
	}
}

func NewIdentityServer(d *Driver) *identityServer {
	return &identityServer{
		Driver: d,
	}
}

func NewVolumeCapabilityAccessMode(mode csi.VolumeCapability_AccessMode_Mode) *csi.VolumeCapability_AccessMode {
	return &csi.VolumeCapability_AccessMode{Mode: mode}
}

func NewControllerServiceCapability(cap csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
	return &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func NewNodeServiceCapability(cap csi.NodeServiceCapability_RPC_Type) *csi.NodeServiceCapability {
	return &csi.NodeServiceCapability{
		Type: &csi.NodeServiceCapability_Rpc{
			Rpc: &csi.NodeServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func RunControllerandNodePublishServer(endpoint string, d *Driver, cs csi.ControllerServer, ns csi.NodeServer) {
	ids := NewIdentityServer(d)

	s := NewNonBlockingGRPCServer()
	s.Start(endpoint, ids, cs, ns)
}

func logGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Infof("GRPC call: %s", info.FullMethod)
	log.Infof("GRPC request: %s", protosanitizer.StripSecrets(req))
	resp, err := handler(ctx, req)
	if err != nil {
		log.Errorf("GRPC error: %v", err)
	} else {
		log.Infof("GRPC response: %s", protosanitizer.StripSecrets(resp))
	}
	return resp, err
}

type VolumeMounter struct {
	path     string
	readOnly bool
}

func (l *VolumeMounter) GetPath() string {
	return l.path
}

func (l *VolumeMounter) GetAttributes() volume.Attributes {
	return volume.Attributes{
		ReadOnly:       l.readOnly,
		Managed:        !l.readOnly,
		SELinuxRelabel: false,
	}
}

func (l *VolumeMounter) CanMount() error {
	return nil
}

func (l *VolumeMounter) SetUp(mounterArgs volume.MounterArgs) error {
	return nil
}

func (l *VolumeMounter) SetUpAt(dir string, mounterArgs volume.MounterArgs) error {
	return nil
}

func (l *VolumeMounter) GetMetrics() (*volume.Metrics, error) {
	return nil, nil
}

func SetVolumeOwnership(path string, readOnly bool, gid string, policy string) error {
	id, err := strconv.Atoi(gid)
	if err != nil {
		return fmt.Errorf("convert %s to int failed with %v", gid, err)
	}
	gidInt64 := int64(id)
	fsGroupChangePolicy := v1.FSGroupChangeOnRootMismatch
	if policy != "" {
		fsGroupChangePolicy = v1.PodFSGroupChangePolicy(policy)
	}
	return volume.SetVolumeOwnership(&VolumeMounter{path: path, readOnly: readOnly}, &gidInt64, &fsGroupChangePolicy, nil)
}
