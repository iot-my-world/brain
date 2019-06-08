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
const CompanyRecordHandlerRetrieve Permission = "CompanyRecordHandler.Retrieve"
const CompanyRecordHandlerCollect Permission = "CompanyRecordHandler.Collect"

const CompanyValidatorValidate Permission = "CompanyValidator.Validate"

const CompanyAdministratorCreate Permission = "CompanyAdministrator.Create"
const CompanyAdministratorUpdateAllowedFields Permission = "CompanyAdministrator.UpdateAllowedFields"

// Client
const ClientRecordHandlerRetrieve Permission = "ClientRecordHandler.Retrieve"
const ClientRecordHandlerCollect Permission = "ClientRecordHandler.Collect"

const ClientValidatorValidate Permission = "ClientValidator.Validate"

const ClientAdministratorUpdateAllowedFields Permission = "ClientAdministrator.UpdateAllowedFields"
const ClientAdministratorCreate Permission = "ClientAdministrator.Create"

// Party
const PartyRegistrarInviteCompanyAdminUser Permission = "PartyRegistrar.InviteCompanyAdminUser"
const PartyRegistrarRegisterCompanyAdminUser Permission = "PartyRegistrar.RegisterCompanyAdminUser"
const PartyRegistrarRegisterCompanyUser Permission = "PartyRegistrar.RegisterCompanyUser"
const PartyRegistrarInviteClientAdminUser Permission = "PartyRegistrar.InviteClientAdminUser"
const PartyRegistrarRegisterClientAdminUser Permission = "PartyRegistrar.RegisterClientAdminUser"
const PartyRegistrarRegisterClientUser Permission = "PartyRegistrar.RegisterClientUser"
const PartyRegistrarInviteUser Permission = "PartyRegistrar.InviteUser"
const PartyRegistrarAreAdminsRegistered Permission = "PartyRegistrar.AreAdminsRegistered"

const PartyAdministratorGetMyParty Permission = "PartyAdministrator.GetMyParty"
const PartyAdministratorRetrieveParty Permission = "PartyAdministrator.RetrieveParty"

const PermissionHandlerGetAllUsersViewPermissions Permission = "PermissionHandler.GetAllUsersViewPermissions"

const TrackingReportLive Permission = "TrackingReport.Live"
const TrackingReportHistorical Permission = "TrackingReport.Historical"

// Barcode Scanner
const BarcodeScannerScan Permission = "BarcodeScanner.Scan"

// SF001 Tracker
const SF001TrackerRecordHandlerCollect Permission = "SF001TrackerRecordHandler.Collect"
