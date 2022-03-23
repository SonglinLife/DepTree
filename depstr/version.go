package depstr


type Version struct{
	Org string
	CleanVersion string
}


func NewVersion(verStr string) Version{
	return Version{
		Org: verStr,
		CleanVersion: verStr,
	}
}