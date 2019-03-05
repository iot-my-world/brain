package criterion

type Type string

// basic criteria
const Text Type = "Text"
const DateRange Type = "DateRange"

// compound criteria
const Or Type = "Or"

// list criteria
const ListText Type = "ListText"

// exact criteria
const ExactText Type = "ExactText"
