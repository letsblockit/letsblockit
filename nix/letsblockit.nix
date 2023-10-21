{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-FViI6r08RsGTYTSKgMJ1XKcuadG5RMua7Ggq2XE79p4=";
  version = "1.0";
}
