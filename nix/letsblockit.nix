{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-z9tEPML/9aDUH8XULDmQTa0PJH0hpI7PfKHV8eHqqJc=";
  version = "1.0";
}
