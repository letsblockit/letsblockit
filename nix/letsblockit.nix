{ buildGoModule, go_1_17, cmd ? "server" }:
buildGoModule.override { go = go_1_17; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-IUvS+kqD6JNGUOnd+19NtDOIXPq+JBSi7ovq1zVo7Hw=";
  version = "1.0";
}
