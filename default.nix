{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "sha256-60k5NfPfZgEbL2THCXuMCU5EXMIMLHHiLQwp1qBJueA=";
  src = ./.;
  doCheck = false;
}
