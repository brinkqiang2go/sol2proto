// Copyright 2018 AMIS Technologies
// This file is part of the sol2proto
//
// The sol2proto is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The sol2proto is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the sol2proto. If not, see <http://www.gnu.org/licenses/>.

package grpc

import (
	"fmt"
	"html/template"
	"io"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/getamis/sirius/util"
)

// Generate a renderable object and required message types from an Ethereum contract ABI
func GenerateServiceProtoFile(srvName, pkgName string, contractABI abi.ABI, version string) (protoFile ProtoFile, msgs []Message) {
	protoFile = ProtoFile{
		GeneratorVersion: version,
		Package:          pkgName,
		Name:             util.ToCamelCase(srvName),
		Sources:          []string{fmt.Sprintf("%s.abi", srvName)},
	}

	methods, requiredMsgs := ParseMethods(contractABI.Methods)
	protoFile.Methods = append(protoFile.Methods, methods...)

	msgs = append(msgs, requiredMsgs...)

	events, requiredMsgs := ParseEvents(contractABI.Events)
	protoFile.Events = append(protoFile.Events, events...)

	msgs = append(msgs, requiredMsgs...)

	sort.Sort(protoFile.Methods)
	sort.Sort(protoFile.Events)
	sort.Sort(protoFile.Sources)

	return protoFile, msgs
}

type ProtoFile struct {
	GeneratorVersion string
	Package          string
	Name             string
	Methods          Methods
	Events           Methods
	Sources          Sources
}

func (p ProtoFile) Render(writer io.WriteCloser) error {
	template, err := template.New("proto").Parse(ServiceTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		return err
	}

	return template.Execute(writer, p)
}

var ServiceTemplate string = `// Automatically generated by sol2proto {{ .GeneratorVersion }}. DO NOT EDIT!
// sources: {{ range .Sources }}
//     {{ . }}
{{- end }}
syntax = "proto3";

package {{ .Package }};

import "messages.proto";

service {{ .Name }} {
{{- range .Methods }}
    {{ . }}
{{- end }}

    // Not supported yet
{{- range .Events }}
    // {{ . }}
{{- end }}
}
`