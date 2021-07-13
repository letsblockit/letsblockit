{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "1q3icwkqfgpjizp7i9mng63m9v8ch90ihs8fqn3lx60cxsf9xllz";
  src = ./.;
}
