package project

import (
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"reflect"
	"testing"
)

func TestServer_ProtoMessage(t *testing.T) {
	s := &Server{}
	if _, ok := interface{}(s).(interface{ ProtoMessage() }); !ok {
		t.Errorf("Server does not implement ProtoMessage()")
	}
}

func TestServer_ProtoReflect(t *testing.T) {
	s := &Server{}
	if reflect.ValueOf(s).MethodByName("ProtoReflect").IsZero() {
		t.Errorf("Server does not implement ProtoReflect()")
	}
}

func TestServer_Descriptor(t *testing.T) {
	s := &Server{}
	desc, _ := interface{}(s).(interface{ Descriptor() ([]byte, []int) })
	if desc == nil {
		t.Errorf("Server does not implement Descriptor()")
	}
	_, ok := interface{}(s).(interface{ Descriptor() ([]byte, []int) })
	if !ok {
		t.Errorf("Server does not implement Descriptor()")
	}
}

func TestServer_GetId(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want string
	}{
		{
			name: "empty id",
			s:    &Server{},
			want: "",
		},
		{
			name: "non-empty id",
			s:    &Server{Id: "test-id"},
			want: "test-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetId()
			if got != tt.want {
				t.Errorf("GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetUserId(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want string
	}{
		{
			name: "empty user_id",
			s:    &Server{},
			want: "",
		},
		{
			name: "non-empty user_id",
			s:    &Server{UserId: "test-user"},
			want: "test-user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetUserId()
			if got != tt.want {
				t.Errorf("GetUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetName(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want string
	}{
		{
			name: "empty name",
			s:    &Server{},
			want: "",
		},
		{
			name: "non-empty name",
			s:    &Server{Name: "test-name"},
			want: "test-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetName()
			if got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetNetworking(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want *Networking
	}{
		{
			name: "nil networking",
			s:    &Server{},
			want: nil,
		},
		{
			name: "non-nil networking",
			s:    &Server{Networking: &Networking{Address: "127.0.0.1", Ports: []int32{80, 443}}},
			want: &Networking{Address: "127.0.0.1", Ports: []int32{80, 443}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetNetworking()
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetNetworking() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestServer_GetStatus(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want int32
	}{
		{
			name: "default status",
			s:    &Server{},
			want: 0,
		},
		{
			name: "custom status",
			s:    &Server{Status: 1},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetStatus()
			if got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetTariffInfo(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want *TariffInfo
	}{
		{
			name: "nil tariff info",
			s:    &Server{},
			want: nil,
		},
		{
			name: "non-nil tariff info",
			s:    &Server{TariffInfo: &TariffInfo{TariffId: 1, TariffStatus: 1, ExpirationTime: "2024-12-31"}},
			want: &TariffInfo{TariffId: 1, TariffStatus: 1, ExpirationTime: "2024-12-31"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetTariffInfo()
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetTariffInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestServer_GetCreatedAt(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want string
	}{
		{
			name: "empty created_at",
			s:    &Server{},
			want: "",
		},
		{
			name: "non-empty created_at",
			s:    &Server{CreatedAt: "2023-11-15"},
			want: "2023-11-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetCreatedAt()
			if got != tt.want {
				t.Errorf("GetCreatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworking_ProtoMessage(t *testing.T) {
	n := &Networking{}
	if _, ok := interface{}(n).(interface{ ProtoMessage() }); !ok {
		t.Errorf("Networking does not implement ProtoMessage()")
	}
}

func TestNetworking_ProtoReflect(t *testing.T) {
	n := &Networking{}
	if reflect.ValueOf(n).MethodByName("ProtoReflect").IsZero() {
		t.Errorf("Networking does not implement ProtoReflect()")
	}
}

func TestNetworking_Descriptor(t *testing.T) {
	n := &Networking{}
	desc, _ := interface{}(n).(interface{ Descriptor() ([]byte, []int) })
	if desc == nil {
		t.Errorf("Networking does not implement Descriptor()")
	}
	_, ok := interface{}(n).(interface{ Descriptor() ([]byte, []int) })
	if !ok {
		t.Errorf("Networking does not implement Descriptor()")
	}
}

func TestNetworking_GetAddress(t *testing.T) {
	tests := []struct {
		name string
		n    *Networking
		want string
	}{
		{
			name: "empty address",
			n:    &Networking{},
			want: "",
		},
		{
			name: "non-empty address",
			n:    &Networking{Address: "127.0.0.1"},
			want: "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.GetAddress()
			if got != tt.want {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworking_GetPorts(t *testing.T) {
	tests := []struct {
		name string
		n    *Networking
		want []int32
	}{
		{
			name: "nil ports",
			n:    &Networking{},
			want: nil,
		},
		{
			name: "non-nil ports",
			n:    &Networking{Ports: []int32{80, 443}},
			want: []int32{80, 443},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.GetPorts()
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetPorts() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTariffInfo_ProtoMessage(t *testing.T) {
	ti := &TariffInfo{}
	if _, ok := interface{}(ti).(interface{ ProtoMessage() }); !ok {
		t.Errorf("TariffInfo does not implement ProtoMessage()")
	}
}

func TestTariffInfo_ProtoReflect(t *testing.T) {
	ti := &TariffInfo{}
	if reflect.ValueOf(ti).MethodByName("ProtoReflect").IsZero() {
		t.Errorf("TariffInfo does not implement ProtoReflect()")
	}
}

func TestTariffInfo_Descriptor(t *testing.T) {
	ti := &TariffInfo{}
	desc, _ := interface{}(ti).(interface{ Descriptor() ([]byte, []int) })
	if desc == nil {
		t.Errorf("TariffInfo does not implement Descriptor()")
	}
	_, ok := interface{}(ti).(interface{ Descriptor() ([]byte, []int) })
	if !ok {
		t.Errorf("TariffInfo does not implement Descriptor()")
	}
}

func TestTariffInfo_GetTariffId(t *testing.T) {
	tests := []struct {
		name string
		ti   *TariffInfo
		want int32
	}{
		{
			name: "default tariff id",
			ti:   &TariffInfo{},
			want: 0,
		},
		{
			name: "custom tariff id",
			ti:   &TariffInfo{TariffId: 1},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ti.GetTariffId()
			if got != tt.want {
				t.Errorf("GetTariffId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTariffInfo_GetTariffStatus(t *testing.T) {
	tests := []struct {
		name string
		ti   *TariffInfo
		want int32
	}{
		{
			name: "default tariff status",
			ti:   &TariffInfo{},
			want: 0,
		},
		{
			name: "custom tariff status",
			ti:   &TariffInfo{TariffStatus: 1},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ti.GetTariffStatus()
			if got != tt.want {
				t.Errorf("GetTariffStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTariffInfo_GetExpirationTime(t *testing.T) {
	tests := []struct {
		name string
		ti   *TariffInfo
		want string
	}{
		{
			name: "empty expiration time",
			ti:   &TariffInfo{},
			want: "",
		},
		{
			name: "non-empty expiration time",
			ti:   &TariffInfo{ExpirationTime: "2024-12-31"},
			want: "2024-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ti.GetExpirationTime()
			if got != tt.want {
				t.Errorf("GetExpirationTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateProjectRequest_String(t *testing.T) {
	tests := []struct {
		name string
		cpr  *CreateProjectRequest
		want string
	}{
		{
			name: "create project request with fields",
			cpr:  &CreateProjectRequest{Name: "test-project"},
			want: `name:"test-project"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cpr.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateProjectRequest_ProtoMessage(t *testing.T) {
	cpr := &CreateProjectRequest{}
	if _, ok := interface{}(cpr).(interface{ ProtoMessage() }); !ok {
		t.Errorf("CreateProjectRequest does not implement ProtoMessage()")
	}
}

func TestCreateProjectRequest_ProtoReflect(t *testing.T) {
	cpr := &CreateProjectRequest{}
	if reflect.ValueOf(cpr).MethodByName("ProtoReflect").IsZero() {
		t.Errorf("CreateProjectRequest does not implement ProtoReflect()")
	}
}

func TestCreateProjectRequest_Descriptor(t *testing.T) {
	cpr := &CreateProjectRequest{}
	desc, _ := interface{}(cpr).(interface{ Descriptor() ([]byte, []int) })
	if desc == nil {
		t.Errorf("CreateProjectRequest does not implement Descriptor()")
	}
	_, ok := interface{}(cpr).(interface{ Descriptor() ([]byte, []int) })
	if !ok {
		t.Errorf("CreateProjectRequest does not implement Descriptor()")
	}
}

func TestCreateProjectRequest_GetName(t *testing.T) {
	tests := []struct {
		name string
		cpr  *CreateProjectRequest
		want string
	}{
		{
			name: "empty name",
			cpr:  &CreateProjectRequest{},
			want: "",
		},
		{
			name: "non-empty name",
			cpr:  &CreateProjectRequest{Name: "test-project"},
			want: "test-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cpr.GetName()
			if got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateProjectResponse_String(t *testing.T) {
	tests := []struct {
		name string
		cpr  *CreateProjectResponse
		want string
	}{
		{
			name: "create project response with fields",
			cpr:  &CreateProjectResponse{Status: true},
			want: `Status:true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cpr.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateProjectResponse_ProtoMessage(t *testing.T) {
	cpr := &CreateProjectResponse{}
	if _, ok := interface{}(cpr).(interface{ ProtoMessage() }); !ok {
		t.Errorf("CreateProjectResponse does not implement ProtoMessage()")
	}
}

func TestCreateProjectResponse_ProtoReflect(t *testing.T) {
	cpr := &CreateProjectResponse{}
	if reflect.ValueOf(cpr).MethodByName("ProtoReflect").IsZero() {
		t.Errorf("CreateProjectResponse does not implement ProtoReflect()")
	}
}

func TestCreateProjectResponse_Descriptor(t *testing.T) {
	cpr := &CreateProjectResponse{}
	desc, _ := interface{}(cpr).(interface{ Descriptor() ([]byte, []int) })
	if desc == nil {
		t.Errorf("CreateProjectResponse does not implement Descriptor()")
	}
	_, ok := interface{}(cpr).(interface{ Descriptor() ([]byte, []int) })
	if !ok {
		t.Errorf("CreateProjectResponse does not implement Descriptor()")
	}
}

func TestCreateProjectResponse_GetStatus(t *testing.T) {
	tests := []struct {
		name string
		cpr  *CreateProjectResponse
		want bool
	}{
		{
			name: "default status",
			cpr:  &CreateProjectResponse{},
			want: false,
		},
		{
			name: "custom status",
			cpr:  &CreateProjectResponse{Status: true},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cpr.GetStatus()
			if got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProjectResponse_ProtoMessage(t *testing.T) {
	gpr := &GetProjectResponse{}
	if _, ok := interface{}(gpr).(interface{ ProtoMessage() }); !ok {
		t.Errorf("GetProjectResponse does not implement ProtoMessage()")
	}
}

func TestGetProjectResponse_ProtoReflect(t *testing.T) {
	gpr := &GetProjectResponse{}
	if reflect.ValueOf(gpr).MethodByName("ProtoReflect").IsZero() {
		t.Errorf("GetProjectResponse does not implement ProtoReflect()")
	}
}

func TestGetProjectResponse_Descriptor(t *testing.T) {
	gpr := &GetProjectResponse{}
	desc, _ := interface{}(gpr).(interface{ Descriptor() ([]byte, []int) })
	if desc == nil {
		t.Errorf("GetProjectResponse does not implement Descriptor()")
	}
	_, ok := interface{}(gpr).(interface{ Descriptor() ([]byte, []int) })
	if !ok {
		t.Errorf("GetProjectResponse does not implement Descriptor()")
	}
}

func TestGetProjectResponse_GetServer(t *testing.T) {
	tests := []struct {
		name string
		gpr  *GetProjectResponse
		want *Server
	}{
		{
			name: "nil server",
			gpr:  &GetProjectResponse{},
			want: nil,
		},
		{
			name: "non-nil server",
			gpr:  &GetProjectResponse{Server: &Server{Id: "1"}},
			want: &Server{Id: "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.gpr.GetServer()
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("GetServer() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
