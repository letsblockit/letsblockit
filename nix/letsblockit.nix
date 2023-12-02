{ buildGoModule, cmd ? "server" }:
buildGoModule {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-tH/NaZTb/1CAOmgAYYfwpcYJ6g6PBHI354QB84aRMQo=";
  version = "1.0";
}
