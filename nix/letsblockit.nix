{ buildGoModule, go_1_19, cmd ? "server" }:
buildGoModule.override { go = go_1_19; } {
  doCheck = false;
  pname = "letsblockit";
  src = ./..;
  subPackages = "cmd/" + cmd;
  vendorSha256 = "sha256-NiVN+STL6QwLl+7sYD1Gs+H69eDDeyOEi2msmrGdUP0=";
  version = "1.0";
}
