package systemRole

import "gitlab.com/iotTracker/brain/security"

var initialRoles = func() []SystemRole {

	// Register roles here
	allRoles := []SystemRole{
		owner,
		admin,
		employee,
	}

	//Register additional root permissions here
	rootPermisions := []security.Permission{
		"Role.Create",
		"Role.Retrieve",
		"Role.Update",
		"Role.Delete",
	}

	// Create root role and apply permissions of all other roles to root
	for _, role := range allRoles {
		rootPermisions = append(rootPermisions, role.Permissions...)
	}
	root := SystemRole{
		Name: "root",
		Permissions: rootPermisions,
	}
	return append([]SystemRole{root}, allRoles...)
}()

// Create Roles here

var owner = SystemRole{
	Name: "owner",
	Permissions: []security.Permission{
		"Employee.Create",
		"Employee.RetrieveAll",
		"Employee.Update",
		"Employee.Delete",
		"Employee.RetrieveByShiftAssignment",
		"Employee.RetrieveByTagID",
		"BusinessRole.Create",
		"BusinessRole.RetrieveAll",
		"BusinessRole.Update",
		"BusinessRole.Delete",
		"Ship.Create",
		"Ship.RetrieveAll",
		"Ship.Update",
		"Ship.Delete",
		"BusinessDayConfig.Create",
		"BusinessDayConfig.Update",
		"BusinessDayConfig.Retrieve",
		// "BusinessDay.Create", // this may never happen directly via service
		"BusinessDay.GetCurrent",
		"BusinessDay.GetBefore",
		"BusinessDay.GetAfter",
		"BusinessDay.UpdateShifts",
		"BusinessDay.GetSelected",
	},
}

var admin = SystemRole{
	Name: "admin",
	Permissions: []security.Permission{
		"Employee.Create",
		"Employee.RetrieveAll",
		"Employee.Update",
	},
}

var employee = SystemRole{
	Name: "employee",
	Permissions: []security.Permission{
	},
}
