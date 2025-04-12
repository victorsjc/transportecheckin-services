package models

type Empresa struct {
    CNPJ                  string `json:"cnpj"`
    RazaoSocial           string `json:"razao_social"`
    EmailRepresentante    string `json:"email_representante"`
    ContatoRepresentante  string `json:"contato_representante"`
}