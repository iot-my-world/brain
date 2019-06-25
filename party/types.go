package party

type Type string

const System Type = "System"
const Client Type = "Client"
const Company Type = "Company"
const Individual Type = "Individual"

func IsValidType(partyType Type) bool {
	allValidTypes := []Type{
		System, Client, Company,
	}

	for _, validType := range allValidTypes {
		if validType == partyType {
			return true
		}
	}
	return false
}
