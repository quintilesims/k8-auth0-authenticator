resource "auth0_client" "this" {
  name            = "${var.name}"
  description     = "Kubernetes Authenticator"
  app_type        = "regular_web"
  is_first_party  = true
  oidc_conformant = true
  callbacks       = ["http://localhost:9090/token", "https://${var.subdomain}.${var.domain}/token"]
  web_origins     = ["http://localhost:9090", "https://${var.subdomain}.${var.domain}"]
  grant_types     = ["implicit", "authorization_code", "refresh_token", "client_credentials"]
}
