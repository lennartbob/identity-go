package core

type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
	MemberRoleSystem MemberRole = "system"
)

func (r MemberRole) IsValid() bool {
	return r == MemberRoleAdmin || r == MemberRoleMember || r == MemberRoleSystem
}

func (r MemberRole) String() string {
	return string(r)
}
