package api

type Permission string

const RoleCreate Permission = "Role.Create"
const RoleRetrieve Permission = "Role.Retrieve"
const RoleUpdate Permission = "Role.Update"
const RoleDelete Permission = "Role.Delete"

// User
const UserRecordHandlerRetrieve Permission = "UserRecordHandler.Retrieve"
const UserRecordHandlerCollect Permission = "UserRecordHandler.Collect"

const UserAdministratorGetMyUser Permission = "UserAdministrator.GetMyUser"
const UserAdministratorCreate Permission = "UserAdministrator.Create"
const UserAdministratorUpdateAllowedFields Permission = "UserAdministrator.UpdateAllowedFields"
const UserAdministratorUpdatePassword Permission = "UserAdministrator.UpdatePassword"
const UserAdministratorCheckPassword Permission = "UserAdministrator.CheckPassword"
const UserAdministratorSetPassword Permission = "UserAdministrator.SetPassword"

const UserValidatorValidate Permission = "UserValidator.Validate"

// API User
const APIUserRecordHandlerCollect Permission = "APIUserRecordHandler.Collect"

const APIUserAdministratorCreate Permission = "APIUserAdministrator.Create"

const APIUserValidatorValidate Permission = "APIUserValidator.Validate"

// System
const SystemRecordHandlerCollect Permission = "SystemRecordHandler.Collect"

const SystemAdministratorUpdateAllowedFields Permission = "SystemAdministrator.UpdateAllowedFields"

// Company

const CompanyValidatorValidate Permission = "CompanyValidator.Validate"

// Party
const PartyRegistrarInviteCompanyAdminUser Permission = "PartyRegistrar.InviteCompanyAdminUser"
const PartyRegistrarRegisterCompanyAdminUser Permission = "PartyRegistrar.RegisterCompanyAdminUser"
const PartyRegistrarRegisterCompanyUser Permission = "PartyRegistrar.RegisterCompanyUser"
const PartyRegistrarInviteClientAdminUser Permission = "PartyRegistrar.InviteClientAdminUser"
const PartyRegistrarRegisterClientAdminUser Permission = "PartyRegistrar.RegisterClientAdminUser"
const PartyRegistrarRegisterClientUser Permission = "PartyRegistrar.RegisterClientUser"
const PartyRegistrarInviteUser Permission = "PartyRegistrar.InviteUser"
const PartyRegistrarAreAdminsRegistered Permission = "PartyRegistrar.AreAdminsRegistered"

const PermissionHandlerGetAllUsersViewPermissions Permission = "PermissionHandler.GetAllUsersViewPermissions"

const TrackingReportLive Permission = "TrackingReport.Live"
const TrackingReportHistorical Permission = "TrackingReport.Historical"

// Barcode Scanner
const BarcodeScannerScan Permission = "BarcodeScanner.Scan"

// SF001 Tracker
const SF001TrackerRecordHandlerCollect Permission = "SF001TrackerRecordHandler.Collect"
const SF001TrackerValidatorValidate Permission = "SF001TrackerValidator.Validate"
const SF001TrackerAdministratorCreate Permission = "SF001TrackerAdministrator.Create"
const SF001TrackerAdministratorUpdateAllowedFields Permission = "SF001TrackerAdministrator.UpdateAllowedFields"
