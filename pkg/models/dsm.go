// Copyright 2021 Synology Inc.

package models

import (
	"bytes"
	"encoding/base64"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/SynologyOpenSource/synology-csi/pkg/utils"
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

	// Templates
	LunNameTemplate             = "{{.Prefix}}-{{.VolumeName}}"
	ShareNameTemplate           = "{{.Prefix}}-{{.VolumeNameCompressedUUIDOnly}}"
	LunDescriptionTemplate      = "{{.PvcNamespace}}:{{.PvcName}}"
	ShareDescriptionTemplate    = "{{.PvcNamespace}}:{{.PvcName}}"
	LunSnapshotNameTemplate     = "{{.Prefix}}-s-{{.VolumeName}}"
	ShareSnapshotNameTemplate   = "{{.Prefix}}-s-{{.VolumeNameCompressedUUIDOnly}}"
	SnapshotDescriptionTemplate = "Snapshot {{.PvcNamespace}}:{{.PvcName}}" // iscsi only

	CompiledLunNameTemplate             *template.Template
	CompiledShareNameTemplate           *template.Template
	CompiledLunDescriptionTemplate      *template.Template
	CompiledShareDescriptionTemplate    *template.Template
	CompiledLunSnapshotNameTemplate     *template.Template
	CompiledShareSnapshotNameTemplate   *template.Template
	CompiledSnapshotDescriptionTemplate *template.Template

	uuidRegex *regexp.Regexp
)

func CompileTemplates() error {
	var err error

	uuidRegex, err = regexp.Compile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}")
	if err != nil {
		return err
	}
	CompiledLunNameTemplate, err = template.New("lun_name").Funcs(sprig.FuncMap()).Parse(LunNameTemplate)
	if err != nil {
		return err
	}
	CompiledShareNameTemplate, err = template.New("share_name").Funcs(sprig.FuncMap()).Parse(ShareNameTemplate)
	if err != nil {
		return err
	}
	CompiledLunDescriptionTemplate, err = template.New("lun_description").Funcs(sprig.FuncMap()).Parse(LunDescriptionTemplate)
	if err != nil {
		return err
	}
	CompiledShareDescriptionTemplate, err = template.New("share_description").Funcs(sprig.FuncMap()).Parse(ShareDescriptionTemplate)
	if err != nil {
		return err
	}
	CompiledLunSnapshotNameTemplate, err = template.New("lun_snapshotName").Funcs(sprig.FuncMap()).Parse(LunSnapshotNameTemplate)
	if err != nil {
		return err
	}
	CompiledShareSnapshotNameTemplate, err = template.New("share_snapshotName").Funcs(sprig.FuncMap()).Parse(ShareSnapshotNameTemplate)
	if err != nil {
		return err
	}
	CompiledSnapshotDescriptionTemplate, err = template.New("lun_snapshotDescription").Funcs(sprig.FuncMap()).Parse(SnapshotDescriptionTemplate)
	if err != nil {
		return err
	}

	log.Infof("lun name template: %s", LunNameTemplate)
	log.Infof("lun description template: %s", LunDescriptionTemplate)
	log.Infof("lun snapshot name template: %s", LunSnapshotNameTemplate)
	log.Infof("snapshot description template: %s", SnapshotDescriptionTemplate)

	log.Infof("share name template: %s", ShareNameTemplate)
	log.Infof("share description template: %s", ShareDescriptionTemplate)
	log.Infof("share snapshot name template: %s", ShareSnapshotNameTemplate)

	return nil
}

type StringGenerator struct {
	VolName                     string
	PvcName                     string
	PvcNamespace                string
	PvName                      string
	NameTemplate                string
	DescriptionTemplate         string
	SnapshotNameTemplate        string
	SnapshotDescriptionTemplate string

	CompiledNameTemplate                *template.Template
	CompiledDescriptionTemplate         *template.Template
	CompiledSnapshotNameTemplate        *template.Template
	CompiledSnapshotDescriptionTemplate *template.Template
}

func NewStringGenerator(volName, protocol string, params map[string]string) (*StringGenerator, error) {
	// prefix, volName, pvcName, pvcNamespace, pvName, lunNameTemplate, lunDescriptionTemplate, snapshotNameTemplate, snapshotDescriptionTemplate string

	const (
		pvcNameKey                     = "csi.storage.k8s.io/pvc/name"
		pvcNamespaceKey                = "csi.storage.k8s.io/pvc/namespace"
		pvNameKey                      = "csi.storage.k8s.io/pv/name"
		nameTemplateKey                = "name_template"
		descriptionTemplateKey         = "description_template"
		snapshotNameTemplateKey        = "snapshot_name_template"
		snapshotDescriptionTemplateKey = "snapshot_description_template"
	)

	pvcName := ""
	pvcNamespace := ""
	pvName := ""
	nameTemplate := ""
	descriptionTemplate := ""
	snapshotNameTemplate := ""
	snapshotDescriptionTemplate := ""

	if params[pvcNameKey] != "" {
		pvcName = params[pvcNameKey]
	}

	if params[pvcNamespaceKey] != "" {
		pvcNamespace = params[pvcNamespaceKey]
	}

	if params[pvNameKey] != "" {
		pvName = params[pvNameKey]
	}

	if params[nameTemplate] != "" {
		nameTemplate = params[nameTemplateKey]
	} else {
		switch protocol {
		case utils.ProtocolIscsi:
			nameTemplate = LunNameTemplate
		case utils.ProtocolSmb:
			nameTemplate = ShareNameTemplate
		}
	}

	if params[descriptionTemplateKey] != "" {
		descriptionTemplate = params[descriptionTemplateKey]
	} else {
		switch protocol {
		case utils.ProtocolIscsi:
			descriptionTemplate = LunDescriptionTemplate
		case utils.ProtocolSmb:
			descriptionTemplate = ShareDescriptionTemplate
		}
	}

	if params[snapshotNameTemplateKey] != "" {
		snapshotNameTemplate = params[snapshotNameTemplateKey]
	} else {
		switch protocol {
		case utils.ProtocolIscsi:
			snapshotNameTemplate = LunSnapshotNameTemplate
		case utils.ProtocolSmb:
			snapshotNameTemplate = ShareSnapshotNameTemplate
		}
	}

	if params[snapshotDescriptionTemplateKey] != "" {
		snapshotDescriptionTemplate = params[snapshotDescriptionTemplateKey]
	} else {
		snapshotDescriptionTemplate = SnapshotDescriptionTemplate
	}

	sg := &StringGenerator{
		VolName:                     volName,
		PvcName:                     pvcName,
		PvcNamespace:                pvcNamespace,
		PvName:                      pvName,
		NameTemplate:                nameTemplate,
		DescriptionTemplate:         descriptionTemplate,
		SnapshotNameTemplate:        snapshotNameTemplate,
		SnapshotDescriptionTemplate: snapshotDescriptionTemplate,
	}

	if err := sg.Compile(); err != nil {
		return nil, err
	}

	return sg, nil
}

func (s *StringGenerator) Compile() error {
	var err error
	s.CompiledNameTemplate, err = template.New("name").Funcs(sprig.FuncMap()).Parse(s.NameTemplate)
	if err != nil {
		return err
	}

	s.CompiledDescriptionTemplate, err = template.New("description").Funcs(sprig.FuncMap()).Parse(s.DescriptionTemplate)
	if err != nil {
		return err
	}

	s.CompiledSnapshotNameTemplate, err = template.New("snapshotName").Funcs(sprig.FuncMap()).Parse(s.SnapshotNameTemplate)
	if err != nil {
		return err
	}

	s.CompiledSnapshotDescriptionTemplate, err = template.New("snapshotDescription").Funcs(sprig.FuncMap()).Parse(s.SnapshotDescriptionTemplate)
	if err != nil {
		return err
	}

	return nil
}

func (s *StringGenerator) GenString(template *template.Template, prefix string) (string, error) {
	output := &bytes.Buffer{}
	pvNameCompressed := s.PvName
	volNameCompressed := s.VolName
	pvNameCompressedUUIDOnly := s.PvName
	volNameCompressedUUIDOnly := s.VolName
	pvNameUUIDOnly := s.PvName
	volNameUUIDOnly := s.VolName

	{
		uuidIndex := uuidRegex.FindStringIndex(s.PvName)
		if uuidIndex != nil {
			parsedUUID, err := uuid.FromString(s.PvName[uuidIndex[0]:uuidIndex[1]])
			if err == nil {
				pvNameCompressed = s.PvName[:uuidIndex[0]] + base64.URLEncoding.EncodeToString(parsedUUID.Bytes())[:22] + s.PvName[uuidIndex[1]:]
				pvNameCompressedUUIDOnly = base64.URLEncoding.EncodeToString(parsedUUID.Bytes())[:22]
				pvNameUUIDOnly = parsedUUID.String()
			}
		}
	}
	{
		uuidIndex := uuidRegex.FindStringIndex(s.VolName)
		if uuidIndex != nil {
			parsedUUID, err := uuid.FromString(s.VolName[uuidIndex[0]:uuidIndex[1]])
			if err == nil {
				volNameCompressed = s.VolName[:uuidIndex[0]] + base64.URLEncoding.EncodeToString(parsedUUID.Bytes())[:22] + s.VolName[uuidIndex[1]:]
				volNameCompressedUUIDOnly = base64.URLEncoding.EncodeToString(parsedUUID.Bytes())[:22]
				volNameUUIDOnly = parsedUUID.String()
			}
		}
	}

	err := template.Execute(output, struct {
		Prefix                       string
		VolumeName                   string
		VolumeNameCompressed         string
		VolumeNameCompressedUUIDOnly string
		VolumeNameUUIDOnly           string
		PvcName                      string
		PvcNamespace                 string
		PvName                       string
		PvNameCompressed             string
		PvNameCompressedUUIDOnly     string
		PvNameUUIDOnly               string
	}{
		Prefix:                       prefix,
		VolumeName:                   s.VolName,
		VolumeNameCompressed:         volNameCompressed,
		VolumeNameCompressedUUIDOnly: volNameCompressedUUIDOnly,
		VolumeNameUUIDOnly:           volNameUUIDOnly,
		PvcName:                      s.PvcName,
		PvcNamespace:                 s.PvcNamespace,
		PvName:                       s.PvName,
		PvNameCompressed:             pvNameCompressed,
		PvNameCompressedUUIDOnly:     pvNameCompressedUUIDOnly,
		PvNameUUIDOnly:               pvNameUUIDOnly,
	})

	if err != nil {
		return "", err
	}
	return output.String(), nil
}

func (s *StringGenerator) GenLunName() (string, error) {
	name, err := s.GenString(s.CompiledNameTemplate, LunPrefix)
	return name, err
}

func (s *StringGenerator) GenShareName() (string, error) {
	name, err := s.GenString(s.CompiledNameTemplate, SharePrefix)

	if len(name) > MaxShareLen {
		name = name[:MaxShareLen]
	}
	return name, err
}

func (s *StringGenerator) GenTargetName() (string, error) {
	name, err := s.GenString(s.CompiledNameTemplate, TargetPrefix)
	return name, err
}

func (s *StringGenerator) GenSnapshotLunName() (string, error) {
	name, err := s.GenString(s.CompiledSnapshotNameTemplate, LunPrefix)
	return name, err
}

func (s *StringGenerator) GenSnapshotShareName() (string, error) {
	name, err := s.GenString(s.CompiledSnapshotNameTemplate, SharePrefix)

	if len(name) > MaxShareLen {
		name = name[:MaxShareLen]
	}
	return name, err
}

func (s *StringGenerator) GenDescription() (string, error) {
	name, err := s.GenString(s.CompiledSnapshotDescriptionTemplate, "")
	return name, err
}

func (s *StringGenerator) GenSnapshotDescription() (string, error) {
	name, err := s.GenString(s.CompiledSnapshotDescriptionTemplate, "")
	return name, err
}
