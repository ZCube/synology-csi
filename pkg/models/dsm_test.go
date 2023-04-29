// Copyright 2021 Synology Inc.

package models

import (
	"fmt"
	"testing"
	"text/template"
)

func TestStringGenerator_GenString(t *testing.T) {
	LunNameTemplate = "{{.VolumeName}}"
	ShareNameTemplate = "{{.VolumeNameCompressed}}"
	err := CompileTemplates()
	if err != nil {
		t.Errorf("GenString() error = %v", err)
		return
	}
	type fields struct {
		VolName                             string
		PvcName                             string
		PvcNamespace                        string
		PvName                              string
		NameTemplate                        string
		DescriptionTemplate                 string
		SnapshotNameTemplate                string
		SnapshotDescriptionTemplate         string
		CompiledNameTemplate                *template.Template
		CompiledDescriptionTemplate         *template.Template
		CompiledSnapshotNameTemplate        *template.Template
		CompiledSnapshotDescriptionTemplate *template.Template
	}
	type args struct {
		template *template.Template
		prefix   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				VolName:      "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				PvcName:      "test-pvc",
				PvcNamespace: "default",
				PvName:       "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				NameTemplate: "{{.VolumeName}}",
			},
			args: args{
				template: CompiledLunNameTemplate,
				prefix:   LunPrefix,
			},
			want: fmt.Sprintf("%v-%v", LunPrefix, "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50"),
		},
		{
			name: "test2",
			fields: fields{
				VolName:      "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				PvcName:      "test-pvc",
				PvcNamespace: "default",
				PvName:       "pvc-9e466146-a2f9-49b6-a8d9-77ba33285c50",
				NameTemplate: "{{.VolumeNameCompressed}}",
			},
			args: args{
				template: CompiledShareNameTemplate,
				prefix:   SharePrefix,
			},
			want: fmt.Sprintf("%v-%v", SharePrefix, "pvc-nkZhRqL5Sbao2Xe6MyhcUA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringGenerator{
				VolName:                             tt.fields.VolName,
				PvcName:                             tt.fields.PvcName,
				PvcNamespace:                        tt.fields.PvcNamespace,
				PvName:                              tt.fields.PvName,
				NameTemplate:                        tt.fields.NameTemplate,
				DescriptionTemplate:                 tt.fields.DescriptionTemplate,
				SnapshotNameTemplate:                tt.fields.SnapshotNameTemplate,
				SnapshotDescriptionTemplate:         tt.fields.SnapshotDescriptionTemplate,
				CompiledNameTemplate:                tt.fields.CompiledNameTemplate,
				CompiledDescriptionTemplate:         tt.fields.CompiledDescriptionTemplate,
				CompiledSnapshotNameTemplate:        tt.fields.CompiledSnapshotNameTemplate,
				CompiledSnapshotDescriptionTemplate: tt.fields.CompiledSnapshotDescriptionTemplate,
			}
			err := s.Compile()
			if err != nil {
				t.Errorf("StringGenerator.GenString() error = %v", err)
				return
			}
			got, err := s.GenString(tt.args.template, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringGenerator.GenString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringGenerator.GenString() = %v, want %v", got, tt.want)
			}
		})
	}
}
