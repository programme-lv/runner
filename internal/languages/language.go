package languages

type ProgrammingLanguage struct {
    Id string `json:"id"`
    FullName string `json:"full_name"`
    CodeFilename string `json:"code_filename"`
    CompileCmd string `json:"compile_cmd"`
    ExecuteCmd string `json:"execute_cmd"`
    EnvVersionCmd string `json:"env_version_cmd"`
    HelloWorldCode string `json:"hello_world_code"`
}
