package script_env

// COUNTRY is embedded func that has access to event and repo
func (e Environment) COUNTRY(lst ...interface{}) (interface{}, error) {
	return e.getFromRepo("country")
}
