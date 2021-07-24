{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "01hw9klkci2w566b14gaircffdaf3z97hai2zf2kvwrbiw9xai3p";
  src = ./.;
}
