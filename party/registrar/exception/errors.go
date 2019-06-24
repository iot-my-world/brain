package exception

import "strings"

type UnableToRetrieveParty struct {
	Reasons []string
}

func (e UnableToRetrieveParty) Error() string {
	return "unable to retrieve party: " + strings.Join(e.Reasons, "; ")
}

type UnableToCollectParties struct {
	Reasons []string
}

func (e UnableToCollectParties) Error() string {
	return "unable to collect parties: " + strings.Join(e.Reasons, "; ")
}

type AlreadyRegistered struct{}

func (e AlreadyRegistered) Error() string {
	return "user already registered"
}

type AlreadyExists struct{}

func (e AlreadyExists) Error() string {
	return "user already exists"
}

type TokenGeneration struct {
	Reasons []string
}

func (e TokenGeneration) Error() string {
	return "token generation error: " + strings.Join(e.Reasons, "; ")
}

type PartyTypeInvalid struct {
	Reasons []string
}

func (e PartyTypeInvalid) Error() string {
	return "party type invalid: " + strings.Join(e.Reasons, "; ")
}

type EmailGeneration struct {
	Reasons []string
}

func (e EmailGeneration) Error() string {
	return "email generation error: " + strings.Join(e.Reasons, "; ")
}

type RegisterSystemAdminUser struct {
	Reasons []string
}

func (e RegisterSystemAdminUser) Error() string {
	return "error registering system admin user: " + strings.Join(e.Reasons, "; ")
}

type InviteCompanyAdminUser struct {
	Reasons []string
}

func (e InviteCompanyAdminUser) Error() string {
	return "error inviting company admin user: " + strings.Join(e.Reasons, "; ")
}

type RegisterCompanyAdminUser struct {
	Reasons []string
}

func (e RegisterCompanyAdminUser) Error() string {
	return "error registering company admin user: " + strings.Join(e.Reasons, "; ")
}

type InviteCompanyUser struct {
	Reasons []string
}

func (e InviteCompanyUser) Error() string {
	return "error inviting company user: " + strings.Join(e.Reasons, "; ")
}

type RegisterCompanyUser struct {
	Reasons []string
}

func (e RegisterCompanyUser) Error() string {
	return "error registering company user: " + strings.Join(e.Reasons, "; ")
}

type InviteClientAdminUser struct {
	Reasons []string
}

func (e InviteClientAdminUser) Error() string {
	return "error inviting client admin user: " + strings.Join(e.Reasons, "; ")
}

type RegisterClientAdminUser struct {
	Reasons []string
}

func (e RegisterClientAdminUser) Error() string {
	return "error registering client admin user: " + strings.Join(e.Reasons, "; ")
}

type InviteClientUser struct {
	Reasons []string
}

func (e InviteClientUser) Error() string {
	return "error inviting client user: " + strings.Join(e.Reasons, "; ")
}
