{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-P7h7+HtVSCab5gfgJbcvgsC8eXF73SUx8sFoAcJkWNE=";
  version = "1.0";
}
