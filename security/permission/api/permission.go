package api

type Permission string

// User

const UserValidatorValidate Permission = "UserValidator.Validate"

// API User

const APIUserAdministratorCreate Permission = "APIUserAdministrator.Create"

const APIUserValidatorValidate Permission = "APIUserValidator.Validate"

// System
const SystemRecordHandlerCollect Permission = "SystemRecordHandler.Collect"

const SystemAdministratorUpdateAllowedFields Permission = "SystemAdministrator.UpdateAllowedFields"

// Company

// Party

const PermissionHandlerGetAllUsersViewPermissions Permission = "PermissionHandler.GetAllUsersViewPermissions"

// Barcode Scanner
const BarcodeScannerScan Permission = "BarcodeScanner.Scan"

// SF001 Tracker
const SF001TrackerRecordHandlerCollect Permission = "SF001TrackerRecordHandler.Collect"
const SF001TrackerValidatorValidate Permission = "SF001TrackerValidator.Validate"
const SF001TrackerAdministratorCreate Permission = "SF001TrackerAdministrator.Create"
const SF001TrackerAdministratorUpdateAllowedFields Permission = "SF001TrackerAdministrator.UpdateAllowedFields"
