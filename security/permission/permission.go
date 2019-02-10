package permission

type Permission string

const RoleCreate = "Role.Create"
const RoleRetrieve = "Role.Retrieve"
const RoleUpdate = "Role.Update"
const RoleDelete = "Role.Delete"

const CompanyRecordHandlerCreate = "CompanyRecordHandler.Create"
const CompanyRecordHandlerRetrieve = "CompanyRecordHandler.Retrieve"
const CompanyRecordHandlerUpdate = "CompanyRecordHandler.Update"
const CompanyRecordHandlerDelete = "CompanyRecordHandler.Delete"
const CompanyRecordHandlerValidate = "CompanyRecordHandler.Validate"
const CompanyRecordHandlerCollect = "CompanyRecordHandler.Collect"
