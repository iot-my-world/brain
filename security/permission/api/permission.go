package api

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
const UserRecordHandlerChangePassword = "UserRecordHandler.ChangePassword"

const CompanyRecordHandlerCreate = "CompanyRecordHandler.Create"
const CompanyRecordHandlerRetrieve = "CompanyRecordHandler.Retrieve"
const CompanyRecordHandlerUpdate = "CompanyRecordHandler.Update"
const CompanyRecordHandlerDelete = "CompanyRecordHandler.Delete"
const CompanyRecordHandlerValidate = "CompanyRecordHandler.Validate"
const CompanyRecordHandlerCollect = "CompanyRecordHandler.Collect"

const ClientRecordHandlerCreate = "ClientRecordHandler.Create"
const ClientRecordHandlerRetrieve = "ClientRecordHandler.Retrieve"
const ClientRecordHandlerUpdate = "ClientRecordHandler.Update"
const ClientRecordHandlerDelete = "ClientRecordHandler.Delete"
const ClientRecordHandlerValidate = "ClientRecordHandler.Validate"
const ClientRecordHandlerCollect = "ClientRecordHandler.Collect"

const PartyRegistrarInviteCompanyAdminUser = "PartyRegistrar.InviteCompanyAdminUser"
const PartyRegistrarRegisterCompanyAdminUser = "PartyRegistrar.RegisterCompanyAdminUser"
const PartyRegistrarInviteClientAdminUser = "PartyRegistrar.InviteClientAdminUser"
const PartyRegistrarRegisterClientAdminUser = "PartyRegistrar.RegisterClientAdminUser"

const PermissionHandlerGetAllUsersViewPermissions = "PermissionHandler.GetAllUsersViewPermissions"

const TK102DeviceRecordHandlerCreate = "TK102DeviceRecordHandler.Create"
const TK102DeviceRecordHandlerRetrieve = "TK102DeviceRecordHandler.Retrieve"
const TK102DeviceRecordHandlerUpdate = "TK102DeviceRecordHandler.Update"
const TK102DeviceRecordHandlerDelete = "TK102DeviceRecordHandler.Delete"
const TK102DeviceRecordHandlerValidate = "TK102DeviceRecordHandler.Validate"
const TK102DeviceRecordHandlerCollect = "TK102DeviceRecordHandler.Collect"

const ReadingRecordHandlerCollect = "ReadingRecordHandler.Collect"
