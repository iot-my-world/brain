package permission

type Permission string

const RoleCreate = "Role.Create"
const RoleRetrieve = "Role.Retrieve"
const RoleUpdate = "Role.Update"
const RoleDelete = "Role.Delete"

const UserRecordHandlerCreate = "UserRecordHandler.Create"
const UserRecordHandlerRetrieve = "UserRecordHandler.Retrieve"
const UserRecordHandlerUpdate = "UserRecordHandler.Update"
const UserRecordHandlerDelete = "UserRecordHandler.Delete"
const UserRecordHandlerValidate = "UserRecordHandler.Validate"
const UserRecordHandlerCollect = "UserRecordHandler.Collect"

const CompanyRecordHandlerCreate = "CompanyRecordHandler.Create"
const CompanyRecordHandlerRetrieve = "CompanyRecordHandler.Retrieve"
const CompanyRecordHandlerUpdate = "CompanyRecordHandler.Update"
const CompanyRecordHandlerDelete = "CompanyRecordHandler.Delete"
const CompanyRecordHandlerValidate = "CompanyRecordHandler.Validate"
const CompanyRecordHandlerCollect = "CompanyRecordHandler.Collect"

const PartyRegistrarInviteCompanyAdminUser = "PartyRegistrar.InviteCompanyAdminUser"
