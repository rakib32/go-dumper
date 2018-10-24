package store

import "github.com/spf13/viper"

type jkey struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientMail              string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderx509CertURL string `json:"auth_provider_x509_cert_url"`
	Clientx509CertURL       string `json:"client_x509_cert_url"`
}

func (j *jkey) loadDef() {
	j.Type = "service_account"
	j.AuthURI = "https://accounts.google.com/o/oauth2/auth"
	j.TokenURI = "https://accounts.google.com/o/oauth2/token"
	j.AuthProviderx509CertURL = "https://www.googleapis.com/oauth2/v1/certss"

}

func (j *jkey) loadConf() {
	j.ProjectID = viper.GetString("src.project_id")
	j.PrivateKeyID = viper.GetString("src.private_key_id")
	j.PrivateKey = viper.GetString("src.private_key")
	j.ClientMail = viper.GetString("src.client_email")
	j.ClientID = viper.GetString("src.client_id")
	j.Clientx509CertURL = viper.GetString("src.certURL")
}
