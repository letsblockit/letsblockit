{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorHash = "sha256-QmISa7nFt6YT67qm1qvPHl+FOH3cI+RtIFrYLREfUfs=";
  version = "1.0";
}
