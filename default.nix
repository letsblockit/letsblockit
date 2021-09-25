{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "1k8i68dnxpxd0q5ndn2rywxffvzngf8pix09iw0q8khs80q9xbrj";
  src = ./.;
}
