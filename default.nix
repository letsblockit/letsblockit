{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "16i98h9xnc4j7dcykk1asvd7iqn5cgjbln7lsxip6r7psajr7f7j";
  src = ./.;
}
