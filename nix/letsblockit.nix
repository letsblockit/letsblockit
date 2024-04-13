{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-mW6BuWw/a9Rf6Wff/QbEdv4zrcrQH9nmIAO9o7DIIAo=";
  version = "1.0";
}
