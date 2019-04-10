package company

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type Company struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	// The email address which will be used to invite the admin
	// user of the company
	// I.e. the first user of the system from the company
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`

	ParentPartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId" bson:"parentId"`
}

type wrapped struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	// The email address which will be used to invite the admin
	// user of the company
	// I.e. the first user of the system from the company
	AdminEmailAddress string `json:"adminEmailAddress" bson:"adminEmailAddress"`

	ParentPartyType party.Type    `json:"parentPartyType" bson:"parentPartyType"`
	ParentId        id.Identifier `json:"parentId" bson:"parentId"`
}

func (c Company) Details() party.Details {
	return party.Details{
		ParentDetail: party.ParentDetail{
			ParentId:        c.ParentId,
			ParentPartyType: c.ParentPartyType,
		},
		Detail: party.Detail{
			PartyId:   id.Identifier{Id: c.Id},
			PartyType: party.Company,
		},
	}
}

func (c Company) SetId(id string) {
	c.Id = id
}

//
//func (c Company) GetBSON() (interface{}, error) {
//	return wrapped{
//		Id:                c.Id,
//		Name:              c.Name,
//		AdminEmailAddress: c.AdminEmailAddress,
//		ParentPartyType:   c.ParentPartyType,
//		ParentId:          c.ParentId,
//	}, nil
//}
//
//func (c Company) SetBSON(raw bson.Raw) error {
//	unmarshalledCompany := new(wrapped)
//	err := raw.Unmarshal(unmarshalledCompany)
//	if err != nil {
//		return err
//	}
//
//	c.Id = unmarshalledCompany.Id
//	c.Name = unmarshalledCompany.Name
//	c.AdminEmailAddress = unmarshalledCompany.AdminEmailAddress
//	c.ParentPartyType = unmarshalledCompany.ParentPartyType
//	c.ParentId = unmarshalledCompany.ParentId
//
//	return nil
//}
