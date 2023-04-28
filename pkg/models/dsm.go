// Copyright 2021 Synology Inc.

package models

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"
	"text/template"

	"github.com/satori/uuid"
	log "github.com/sirupsen/logrus"
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
	TargetPrefix = "csi"
	LunPrefix    = "csi"
	IqnPrefix    = "iqn.2000-01.com.synology:"
	SharePrefix  = "csi"

	// Template
	LunNameTemplate             = "{{.Prefix}}-{{.VolumeName}}"
	ShareNameTemplate           = "{{.Prefix}}-{{.VolumeNameCompressed}}"
	LunDescriptionTemplate      = "{{.PvcNamespace}}:{{.PvcName}}"
	ShareDescriptionTemplate    = "{{.PvcNamespace}}:{{.PvcName}}"
	SnapshotNameTemplate        = "{{.Prefix}}-s-{{.VolumeName}}"
	SnapshotDescriptionTemplate = "Snapshot {{.PvcNamespace}}:{{.PvcName}}"

	CompiledLunNameTemplate             *template.Template
	CompiledShareNameTemplate           *template.Template
	CompiledLunDescriptionTemplate      *template.Template
	CompiledShareDescriptionTemplate    *template.Template
	CompiledSnapshotNameTemplate        *template.Template
	CompiledSnapshotDescriptionTemplate *template.Template

	uuidRegex *regexp.Regexp
)

func CompileTemplates() error {
	var err error

	uuidRegex, err = regexp.Compile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}")
	if err != nil {
		return err
	}
	CompiledLunNameTemplate, err = template.New("name").Parse(LunNameTemplate)
	if err != nil {
		return err
	}
	CompiledShareNameTemplate, err = template.New("name").Parse(ShareNameTemplate)
	if err != nil {
		return err
	}
	CompiledLunDescriptionTemplate, err = template.New("lun_description").Parse(LunDescriptionTemplate)
	if err != nil {
		return err
	}
	CompiledShareDescriptionTemplate, err = template.New("share_description").Parse(ShareDescriptionTemplate)
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

	log.Infof("lun name template: %s", LunNameTemplate)
	log.Infof("share name template: %s", ShareNameTemplate)
	log.Infof("lun description template: %s", LunDescriptionTemplate)
	log.Infof("share description template: %s", ShareDescriptionTemplate)
	log.Infof("snapshot name template: %s", SnapshotNameTemplate)
	log.Infof("snapshot description template: %s", SnapshotDescriptionTemplate)

	return nil
}

type StringGenerator struct {
	VolName      string
	PvcName      string
	PvcNamespace string
	PvName       string
}

func GenLunName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledLunNameTemplate, LunPrefix, volName, pvcName, pvcNamespace, pvName, "")
	return name
}

func GenShareName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledShareNameTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName, "")

	if len(name) > MaxShareLen {
		return name[:MaxShareLen]
	}
	return name
}

func GenTargetName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledLunNameTemplate, TargetPrefix, volName, pvcName, pvcNamespace, pvName, "")
	return name
}

func GenSnapshotLunName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledSnapshotNameTemplate, LunPrefix, volName, pvcName, pvcNamespace, pvName, "")
	return name
}

func GenSnapshotShareName(volName, pvcName, pvcNamespace, pvName string) string {
	name := GenString(CompiledSnapshotNameTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName, "")

	if len(name) > MaxShareLen {
		return name[:MaxShareLen]
	}
	return name
}

func GenLunDescription(volName, pvcName, pvcNamespace, pvName string) string {
	desc := GenString(CompiledLunDescriptionTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName, "")
	return desc
}

func GenShareDescription(volName, pvcName, pvcNamespace, pvName string) string {
	desc := GenString(CompiledShareDescriptionTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName, "")
	return desc
}

func GenSnapshotDescription(volName, pvcName, pvcNamespace, pvName, description string) string {
	desc := GenString(CompiledSnapshotDescriptionTemplate, SharePrefix, volName, pvcName, pvcNamespace, pvName, description)
	return desc
}

func GenString(template *template.Template, prefix, volName, pvcName, pvcNamespace, pvName, description string) string {
	output := &bytes.Buffer{}
	pvNameCompressed := pvName
	volNameCompressed := volName
	{
		{
			uuidIndex := uuidRegex.FindStringIndex(pvName)
			if uuidIndex != nil {
				parsedUUID, err := uuid.FromString(pvName[uuidIndex[0]:uuidIndex[1]])
				if err == nil {
					pvNameCompressed = pvName[:uuidIndex[0]] + base64.URLEncoding.EncodeToString(parsedUUID.Bytes())[:22] + pvName[uuidIndex[1]:]
				}
			}
		}
		{
			uuidIndex := uuidRegex.FindStringIndex(volName)
			if uuidIndex != nil {
				parsedUUID, err := uuid.FromString(volName[uuidIndex[0]:uuidIndex[1]])
				if err == nil {
					volNameCompressed = pvName[:uuidIndex[0]] + base64.URLEncoding.EncodeToString(parsedUUID.Bytes())[:22] + pvName[uuidIndex[1]:]
				}
			}
		}
	}

	err := template.Execute(output, struct {
		Prefix               string
		VolumeName           string
		VolumeNameCompressed string
		PvcName              string
		PvcNamespace         string
		PvName               string
		PvNameCompressed     string
		Description          string
	}{
		Prefix:               prefix,
		VolumeName:           volName,
		VolumeNameCompressed: volNameCompressed,
		PvcName:              pvcName,
		PvcNamespace:         pvcNamespace,
		PvName:               pvName,
		PvNameCompressed:     pvNameCompressed,
		Description:          description,
	})
	outputString := ""
	if err != nil {
		outputString = fmt.Sprintf("%s-%s", prefix, volName)
	} else {
		outputString = output.String()
	}
	return outputString
}
