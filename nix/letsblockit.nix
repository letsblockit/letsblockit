{ buildGoModule, go_1_18, cmd ? "server" }:
buildGoModule.override { go = go_1_18; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-7o3zfphjqnmTqWSFdPWDFPoc7ezZxEKgj09GFD45XS0=";
  version = "1.0";
}
