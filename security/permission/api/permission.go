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

// TK102 Tracking Device
const TK102DeviceRecordHandlerRetrieve Permission = "TK102DeviceRecordHandler.Retrieve"
const TK102DeviceRecordHandlerCollect Permission = "TK102DeviceRecordHandler.Collect"

const TK102DeviceValidatorValidate Permission = "TK102DeviceValidator.Validate"

const TK102DeviceAdministratorChangeOwnershipAndAssignment Permission = "TK102DeviceAdministrator.ChangeOwnershipAndAssignment"
const TK102DeviceAdministratorCreate Permission = "TK102DeviceAdministrator.Create"

const TK102ReadingRecordHandlerCollect Permission = "TK102ReadingRecordHandler.Collect"
const TK102ReadingAdministratorCreate Permission = "TK102ReadingAdministrator.Create"
const TK102ReadingAdministratorCreateBulk Permission = "TK102ReadingAdministrator.CreateBulk"

// ZX303 Tracking Device
const ZX303DeviceRecordHandlerRetrieve Permission = "ZX303DeviceRecordHandler.Retrieve"
const ZX303DeviceRecordHandlerCollect Permission = "ZX303DeviceRecordHandler.Collect"

const ZX303DeviceValidatorValidate Permission = "ZX303DeviceValidator.Validate"

const ZX303DeviceAdministratorChangeOwnershipAndAssignment Permission = "ZX303DeviceAdministrator.ChangeOwnershipAndAssignment"
const ZX303DeviceAdministratorCreate Permission = "ZX303DeviceAdministrator.Create"
const ZX303DeviceAdministratorUpdateAllowedFields Permission = "ZX303DeviceAdministrator.UpdateAllowedFields"
const ZX303DeviceAdministratorHeartbeat Permission = "ZX303DeviceAdministrator.Heartbeat"

const ZX303DeviceAuthenticatorLogin Permission = "ZX303DeviceAuthenticator.Login"
const ZX303DeviceAuthenticatorLogout Permission = "ZX303DeviceAuthenticator.Logout"

const ZX303TaskRecordHandlerCollect Permission = "ZX303TaskRecordHandler.Collect"
const ZX303TaskValidatorValidate Permission = "ZX303TaskValidator.Validate"
const ZX303TaskAdministratorSubmit Permission = "ZX303TaskAdministrator.Submit"
const ZX303TaskAdministratorFailTask Permission = "ZX303TaskAdministrator.FailTask"
const ZX303TaskAdministratorTransitionTask Permission = "ZX303TaskAdministrator.TransitionTask"

const ZX303StatusReadingReportGeneratorBatteryReport Permission = "ZX303StatusReadingReportGenerator.BatteryReport"

const TrackingReportLive Permission = "TrackingReport.Live"
const TrackingReportHistorical Permission = "TrackingReport.Historical"

// Barcode Scanner
const BarcodeScannerScan Permission = "BarcodeScanner.Scan"
