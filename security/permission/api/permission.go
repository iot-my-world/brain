package api

type Permission string

const RoleCreate Permission = "Role.Create"
const RoleRetrieve Permission = "Role.Retrieve"
const RoleUpdate Permission = "Role.Update"
const RoleDelete Permission = "Role.Delete"

const UserRecordHandlerCreate Permission = "UserRecordHandler.Create"
const UserRecordHandlerRetrieve Permission = "UserRecordHandler.Retrieve"
const UserRecordHandlerUpdate Permission = "UserRecordHandler.Update"
const UserRecordHandlerDelete Permission = "UserRecordHandler.Delete"
const UserRecordHandlerValidate Permission = "UserRecordHandler.Validate"
const UserRecordHandlerCollect Permission = "UserRecordHandler.Collect"
const UserRecordHandlerChangePassword Permission = "UserRecordHandler.ChangePassword"

const SystemRecordHandlerCollect Permission = "SystemRecordHandler.Collect"

const CompanyRecordHandlerCreate Permission = "CompanyRecordHandler.Create"
const CompanyRecordHandlerRetrieve Permission = "CompanyRecordHandler.Retrieve"
const CompanyRecordHandlerUpdate Permission = "CompanyRecordHandler.Update"
const CompanyRecordHandlerDelete Permission = "CompanyRecordHandler.Delete"
const CompanyRecordHandlerValidate Permission = "CompanyRecordHandler.Validate"
const CompanyRecordHandlerCollect Permission = "CompanyRecordHandler.Collect"

const ClientRecordHandlerCreate Permission = "ClientRecordHandler.Create"
const ClientRecordHandlerRetrieve Permission = "ClientRecordHandler.Retrieve"
const ClientRecordHandlerUpdate Permission = "ClientRecordHandler.Update"
const ClientRecordHandlerDelete Permission = "ClientRecordHandler.Delete"
const ClientRecordHandlerValidate Permission = "ClientRecordHandler.Validate"
const ClientRecordHandlerCollect Permission = "ClientRecordHandler.Collect"

const PartyRegistrarInviteCompanyAdminUser Permission = "PartyRegistrar.InviteCompanyAdminUser"
const PartyRegistrarRegisterCompanyAdminUser Permission = "PartyRegistrar.RegisterCompanyAdminUser"
const PartyRegistrarInviteCompanyUser Permission = "PartyRegistrar.InviteCompanyUser"
const PartyRegistrarRegisterCompanyUser Permission = "PartyRegistrar.RegisterCompanyUser"
const PartyRegistrarInviteClientAdminUser Permission = "PartyRegistrar.InviteClientAdminUser"
const PartyRegistrarRegisterClientAdminUser Permission = "PartyRegistrar.RegisterClientAdminUser"
const PartyRegistrarInviteClientUser Permission = "PartyRegistrar.InviteClientUser"
const PartyRegistrarRegisterClientUser Permission = "PartyRegistrar.RegisterClientUser"
const PartyRegistrarAreAdminsRegistered Permission = "PartyRegistrar.AreAdminsRegistered"

const PartyHandlerGetMyParty Permission = "PartyHandler.GetMyParty"

const PermissionHandlerGetAllUsersViewPermissions Permission = "PermissionHandler.GetAllUsersViewPermissions"

const TK102DeviceRecordHandlerCreate Permission = "TK102DeviceRecordHandler.Create"
const TK102DeviceRecordHandlerRetrieve Permission = "TK102DeviceRecordHandler.Retrieve"
const TK102DeviceRecordHandlerDelete Permission = "TK102DeviceRecordHandler.Delete"
const TK102DeviceRecordHandlerValidate Permission = "TK102DeviceRecordHandler.Validate"
const TK102DeviceRecordHandlerCollect Permission = "TK102DeviceRecordHandler.Collect"

const TK102DeviceAdministratorChangeOwner Permission = "TK102DeviceAdministrator.ChangeOwner"
const TK102DeviceAdministratorChangeAssigned Permission = "TK102DeviceAdministrator.ChangeAssigned"

const ReadingRecordHandlerCollect Permission = "ReadingRecordHandler.Collect"

const TrackingReportLive Permission = "TrackingReport.Live"
const TrackingReportHistorical Permission = "TrackingReport.Historical"
