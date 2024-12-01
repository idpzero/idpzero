---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "idpzero"
  text: "Single binary Identity Provider (IdP)"
  tagline: Built for developers.
  actions:
    - theme: brand
      text: Get Started Now
      link: /guide/getting-started
    - theme: alt
      text: Wait. What is it?
      link: /guide/what-is-it

features:
  - title: Single CLI for simplicity
    details: Configure and serve an OpenID Connect and OAuth2 IDP with a single CLI binary.
  - title: Commit to Git 
    details: Configuration is designed to be added to your repo and shared across teams.
  - title: Delightful Login Experience
    details: During login flows, select the user to login as, no need passwords needed.
---

### Simple Login Experience

Simplify your dev/test experience though the simplified OpenID Connect / OAuth2 login experience <span class="idpzero-text">idpzero</span> provides. As part of the login flow, simply select the user to login as from the drop down list. No passwords to worry about. Users and the claims that they have are managed in your local `.idpzero` configuration directory.

![Login](/screenshots/login.png)