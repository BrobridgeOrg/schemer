package schemer

type Runtime interface {
	SetEnv(value map[string]interface{})
	LoadScript(script string) error
	Compile(script string) error
	Execute(sourceSchema *Schema, data map[string]interface{}) ([]map[string]interface{}, error)
}
