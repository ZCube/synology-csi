// Copyright 2021 Synology Inc.

package models

import (
	"bytes"
	"fmt"
	"text/template"
)

const (
	K8sCsiName = "Kubernetes CSI"

	// ISCSI definitions
	FsTypeExt4       = "ext4"
	FsTypeBtrfs      = "btrfs"
	LunTypeFile      = "FILE"
	LunTypeThin      = "THIN"
	LunTypeAdv       = "ADV"
	LunTypeBlun      = "BLUN"       // thin provision, mapped to type 263
	LunTypeBlunThick = "BLUN_THICK" // thick provision, mapped to type 259
	MaxIqnLen        = 128
	MaxLunDescLen    = 127

	// Share definitions
	MaxShareLen             = 32
	MaxShareDescLen         = 64
	UserGroupTypeLocalUser  = "local_user"
	UserGroupTypeLocalGroup = "local_group"
	UserGroupTypeSystem     = "system"

	// CSI definitions
	ShareSnapshotDescPrefix = "(Do not change)"
)

var (
	TargetPrefix = "k8s-csi"
	LunPrefix    = "k8s-csi"
	IqnPrefix    = "iqn.2000-01.com.synology:"
	SharePrefix  = "k8s-csi"

	// Template
	NameTemplate                = "{{.Prefix}}-{{.VolumeName}}"
	DescriptionTemplate         = "{{.PvcNamespace}}:{{.PvcName}}"
	SnapshotNameTemplate        = "{{.Prefix}}-s-{{.VolumeName}}"
	SnapshotDescriptionTemplate = "Snapshot {{.PvcNamespace}}:{{.PvcName}}"

	CompiledNameTemplate                *template.Template
	CompiledDescriptionTemplate         *template.Template
	CompiledSnapshotNameTemplate        *template.Template
	CompiledSnapshotDescriptionTemplate *template.Template
)

func CompileTemplates() error {
	var err error
	CompiledNameTemplate, err = template.New("name").Parse(NameTemplate)
	if err != nil {
		return err
	}
	CompiledDescriptionTemplate, err = template.New("description").Parse(DescriptionTemplate)
	if err != nil {
		return err
	}
	CompiledSnapshotNameTemplate, err = template.New("snapshotName").Parse(SnapshotNameTemplate)
	if err != nil {
		return err
	}
	CompiledSnapshotDescriptionTemplate, err = template.New("snapshotDescription").Parse(SnapshotDescriptionTemplate)
	if err != nil {
		return err
	}
	return nil
}

func GenLunName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledNameTemplate, LunPrefix, volName, pvcName, pvcNamespace, pvName)
	return name
}

func GenShareName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledNameTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName)

	if len(name) > MaxShareLen {
		return name[:MaxShareLen]
	}
	return name
}

func GenTargetName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledNameTemplate, TargetPrefix, volName, pvcName, pvcNamespace, pvName)
	return name
}

func GenSnapshotLunName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledSnapshotNameTemplate, LunPrefix, volName, pvcName, pvcNamespace, pvName)
	return name
}

func GenSnapshotShareName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledSnapshotNameTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName)

	if len(name) > MaxShareLen {
		return name[:MaxShareLen]
	}
	return name
}

func GenDescription(volName, pvcName, pvcNamespace, pvName string) string {
	desc := GenString(CompiledDescriptionTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName)
	return desc
}

func GenSnapshotDescription(volName, pvcName, pvcNamespace, pvName string) string {
	desc := GenString(CompiledSnapshotDescriptionTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName)
	return desc
}

func GenString(template *template.Template, prefix, volName, pvcName, pvcNamespace, pvName string) string {
	output := &bytes.Buffer{}
	err := template.ExecuteTemplate(output, "name", struct {
		Prefix       string
		VolumeName   string
		PvcName      string
		PvcNamespace string
		PvName       string
	}{
		Prefix:       prefix,
		VolumeName:   volName,
		PvcName:      pvcName,
		PvcNamespace: pvcNamespace,
		PvName:       pvName,
	})
	outputString := ""
	if err != nil {
		outputString = fmt.Sprintf("%s-%s", prefix, volName)
	} else {
		outputString = output.String()
	}
	return outputString
}
