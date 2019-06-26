package api

type Permission string

// User

// API User

// System

const SystemAdministratorUpdateAllowedFields Permission = "SystemAdministrator.UpdateAllowedFields"

// Company

// Party

// SF001 Tracker
const SF001TrackerValidatorValidate Permission = "SF001TrackerValidator.Validate"
const SF001TrackerAdministratorCreate Permission = "SF001TrackerAdministrator.Create"
const SF001TrackerAdministratorUpdateAllowedFields Permission = "SF001TrackerAdministrator.UpdateAllowedFields"
