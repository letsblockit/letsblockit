{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "0mkr4v3ps9js8zj5w9rzwlndskwxd5ikv78lx0bba96ssafnkfbs";
  src = ./.;
}
