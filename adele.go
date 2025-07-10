package adele

const Version = "v0.0.0"

// Create a new instance of the Adele type using a pointer to Adele with the
// root path of the application as a argument. The new-up is called by project adele's consuming package
// to bootstrap the framework.
func (a *Adele) New(rootPath string) error {

	// ...

	return nil
}
