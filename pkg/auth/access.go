package auth

type Aces struct {
	Admin int
	Audit int
	Setup int
}

const (
	Admin = 2
	Audit = 3
	Setup = 5
)

func GetAces(stauts int) Aces {
	var (
		aces Aces
	)
	if stauts%Admin == 0 {
		aces.Admin = 1
	}
	if stauts%Audit == 0 {
		aces.Audit = 1
	}
	if stauts%Setup == 0 {
		aces.Setup = 1
	}
	return aces
}
