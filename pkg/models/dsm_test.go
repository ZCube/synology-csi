// Copyright 2021 Synology Inc.

package models

import (
	"fmt"
	"testing"
	"text/template"
)

func TestGenString(t *testing.T) {
	LunNameTemplate = "{{.Prefix}}-{{.VolumeName}}"
	ShareNameTemplate = "{{.Prefix}}-{{.VolumeNameCompressed}}"
	err := CompileTemplates()
	if err != nil {
		t.Errorf("GenString() error = %v", err)
		return
	}
	type args struct {
		template     *template.Template
		prefix       string
		volName      string
		pvcName      string
		pvcNamespace string
		pvName       string
		description  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				template:     CompiledLunNameTemplate,
				prefix:       LunPrefix,
				volName:      "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				pvcName:      "test-pvc",
				pvcNamespace: "default",
				pvName:       "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				description:  "",
			},
			want: fmt.Sprintf("%v-%v", LunPrefix, "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50"),
		},
		{
			name: "test2",
			args: args{
				template:     CompiledShareNameTemplate,
				prefix:       SharePrefix,
				volName:      "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				pvcName:      "test-pvc",
				pvcNamespace: "default",
				pvName:       "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				description:  "",
			},
			want: fmt.Sprintf("%v-%v", SharePrefix, "pvc-nkZhRqL5Sbao2Xe6MyhcUA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenString(tt.args.template, tt.args.prefix, tt.args.volName, tt.args.pvcName, tt.args.pvcNamespace, tt.args.pvName, tt.args.description); got != tt.want {
				t.Errorf("GenString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenShareName(t *testing.T) {
	LunNameTemplate = "{{.Prefix}}-{{.VolumeName}}"
	ShareNameTemplate = "{{.Prefix}}-{{.VolumeNameCompressed}}"
	err := CompileTemplates()
	if err != nil {
		t.Errorf("GenString() error = %v", err)
		return
	}
	type args struct {
		volName      string
		pvcName      string
		pvcNamespace string
		pvName       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test2",
			args: args{
				volName:      "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				pvcName:      "test-pvc",
				pvcNamespace: "default",
				pvName:       "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
			},
			want: fmt.Sprintf("%v-%v", SharePrefix, "pvc-nkZhRqL5Sbao2Xe6MyhcUA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenShareName(tt.args.volName, tt.args.pvcName, tt.args.pvcNamespace, tt.args.pvName); got != tt.want {
				t.Errorf("GenShareName() = %v, want %v", got, tt.want)
			}
		})
	}
}
