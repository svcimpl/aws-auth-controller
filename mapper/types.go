package mapper

import (
	"fmt"
	"strings"
)

//RoleAuthMap is the basic structure of mapRoles authentication object
type RoleAuthMap struct {
	RoleARN  string   `yaml:"rolearn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups,omitempty"`
}

func (r *RoleAuthMap) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("- rolearn: %v\n", r.RoleARN))
	s.WriteString(fmt.Sprintf("  username: %v\n", r.Username))
	s.WriteString(fmt.Sprintf("  groups:\n"))
	for _, group := range r.Groups {
		s.WriteString(fmt.Sprintf("   - %v\n", group))
	}
	return s.String()
}

//UsersAuthMap is the basic structure of a mapUsers authentication object
type UsersAuthMap struct {
	UserARN  string   `yaml:"userarn"`
	UserName string   `yaml:"username"`
	Groups   []string `yaml:"groups,omitempty"`
}

func (u *UsersAuthMap) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("- userarn: %v\n", u.UserARN))
	s.WriteString(fmt.Sprintf("  username: %v\n", u.UserName))
	s.WriteString(fmt.Sprintf("  groups:\n"))
	for _, group := range u.Groups {
		s.WriteString(fmt.Sprintf("   - %v\n", group))
	}
	return s.String()
}

//AWSAuthData is the Data portion of the aws_auth config map
type AWSAuthData struct {
	MapRoles []*RoleAuthMap  `yaml:"mapRoles"`
	MapUsers []*UsersAuthMap `yaml:"mapUsers"`
}
