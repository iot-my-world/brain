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

const DeviceRecordHandlerCreate = "DeviceRecordHandler.Create"
const DeviceRecordHandlerRetrieve = "DeviceRecordHandler.Retrieve"
const DeviceRecordHandlerUpdate = "DeviceRecordHandler.Update"
const DeviceRecordHandlerDelete = "DeviceRecordHandler.Delete"
const DeviceRecordHandlerValidate = "DeviceRecordHandler.Validate"
const DeviceRecordHandlerCollect = "DeviceRecordHandler.Collect"

const ReadingRecordHandlerCollect = "ReadingRecordHandler.Collect"