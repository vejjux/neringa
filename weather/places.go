package weather

type Place struct {
	Code string
	Name string
}

var (
	Nida       = Place{Code: "neringa-nida", Name: "Nida"}
	Preila     = Place{Code: "neringa-preila", Name: "Preila"}
	Pervalka   = Place{Code: "neringa-pervalka", Name: "Pervalka"}
	Juodkrante = Place{Code: "neringa-juodkrante", Name: "Juodkrantė"}
	Smiltyne   = Place{Code: "klaipeda-smiltyne", Name: "Smiltynė"}
)
