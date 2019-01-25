package user

type RecordHandler interface {
	Create(request *CreateRequest, response *CreateResponse) error
	RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error
	Update(request *UpdateRequest, response *UpdateResponse) error
	Delete(request *DeleteRequest, response *DeleteResponse) error
}

type CreateRequest struct{
	NewUser NewUser `json:"newUser"`
}

type CreateResponse struct {
	User User `json:"user"`
}

type RetrieveAllRequest struct {

}

type RetrieveAllResponse struct {
	UserRecords []User `json:"userRecords" bson:"userRecords"`
}

type DeleteRequest struct{
	Username string `json:"username" bson:"username"`
}

type DeleteResponse struct {
}

type UpdateRequest struct{
	UpdatedUser User `json:"updatedUser"`
}

type UpdateResponse struct{
	User User `json:"user"`
}

//type RetrieveRequest struct {
//	Username string `json:"username" bson:"username"`
//}
//
//type RetrieveResponse struct {
//	Success bool `json:"success" bson:"success"`
//	Reasons []string `json:"reasons" bson:"reasons"`
//	User User `json:"user" bson:"user"`
//}
//
//type UpdateRequest struct {
//	Username string `json:"username" bson:"username"`
//	Password string `json:"pwd" bson:"pwd"`
//	SystemRole string `json:"systemRole" bson:"systemRole"`
//	BusinessRole string `json:"businessRole" bson:"businessRole"`
//
//}
//
//type UpdateResponse struct {
//	Success bool `json:"success" bson:"success"`
//	Reasons []string `json:"reasons" bson:"reasons"`
//}
//

