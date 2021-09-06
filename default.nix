{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "0bngv079wjm8x0k78fqcl0ck0k19lbjcyqqc5crxpad4pyry1jdf";
  src = ./.;
}
